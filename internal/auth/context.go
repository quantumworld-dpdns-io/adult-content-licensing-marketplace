package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

var ErrUnauthorized = errors.New("unauthorized")

type Principal struct {
	Role string
	Tier string
	Sub  string
}

func PrincipalFromRequest(r *http.Request) (Principal, error) {
	if p, err := principalFromBearer(r); err == nil {
		return p, nil
	}
	return principalFromHeaders(r)
}

func principalFromBearer(r *http.Request) (Principal, error) {
	a := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(strings.ToLower(a), "bearer ") {
		return Principal{}, ErrUnauthorized
	}
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	token := strings.TrimSpace(a[7:])
	claims, err := VerifyToken(secret, token)
	if err != nil {
		return Principal{}, ErrUnauthorized
	}
	return principalFromRoleTier(claims.Role, claims.Tier, claims.Sub)
}

func principalFromHeaders(r *http.Request) (Principal, error) {
	role := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Role")))
	tier := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Tier")))
	sub := strings.TrimSpace(r.Header.Get("X-Sub"))
	return principalFromRoleTier(role, tier, sub)
}

func principalFromRoleTier(role, tier, sub string) (Principal, error) {
	if role == "" {
		return Principal{}, ErrUnauthorized
	}
	p := Principal{Role: role, Tier: tier, Sub: sub}
	switch role {
	case "admin", "auditor":
		return p, nil
	case "regularuser":
		switch tier {
		case "alpha", "beta", "vip1", "vip2", "vip3":
			return p, nil
		default:
			return Principal{}, ErrUnauthorized
		}
	default:
		return Principal{}, ErrUnauthorized
	}
}

func CanListLicenses(p Principal) bool {
	return p.Role == "admin" || p.Role == "auditor" || p.Role == "regularuser"
}

func CanCreateLicense(p Principal) bool {
	if p.Role == "admin" {
		return true
	}
	if p.Role != "regularuser" {
		return false
	}
	return p.Tier == "vip1" || p.Tier == "vip2" || p.Tier == "vip3"
}
