package auth

import (
	"testing"
	"time"
)

func TestSignAndVerifyToken(t *testing.T) {
	secret := "s3cr3t"
	tok, err := SignToken(secret, Claims{Sub: "u1", Role: "regularuser", Tier: "vip1", Exp: time.Now().Add(time.Hour).Unix()})
	if err != nil {
		t.Fatalf("sign err: %v", err)
	}
	claims, err := VerifyToken(secret, tok)
	if err != nil {
		t.Fatalf("verify err: %v", err)
	}
	if claims.Role != "regularuser" || claims.Tier != "vip1" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
