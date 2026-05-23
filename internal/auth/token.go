package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Claims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	Tier string `json:"tier,omitempty"`
	Exp  int64  `json:"exp"`
}

func SignToken(secret string, claims Claims) (string, error) {
	if secret == "" {
		return "", errors.New("secret is required")
	}
	if claims.Exp == 0 {
		claims.Exp = time.Now().Add(12 * time.Hour).Unix()
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}
	enc := base64.RawURLEncoding.EncodeToString(payload)
	sig := sign(secret, enc)
	return enc + "." + sig, nil
}

func VerifyToken(secret, token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, ErrInvalidToken
	}
	payloadEnc := parts[0]
	sig := parts[1]
	if !hmac.Equal([]byte(sign(secret, payloadEnc)), []byte(sig)) {
		return Claims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(payloadEnc)
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var c Claims
	if err := json.Unmarshal(payload, &c); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if c.Exp < time.Now().Unix() {
		return Claims{}, ErrExpiredToken
	}
	return c, nil
}

func sign(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
