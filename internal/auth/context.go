package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrUnauthorized = errors.New("unauthorized")

type Principal struct {
	Role string
	Tier string
}

func PrincipalFromRequest(r *http.Request) (Principal, error) {
	role := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Role")))
	tier := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Tier")))

	if role == "" {
		return Principal{}, ErrUnauthorized
	}

	switch role {
	case "admin", "auditor":
		return Principal{Role: role, Tier: tier}, nil
	case "regularuser":
		if tier == "" {
			return Principal{}, ErrUnauthorized
		}
		switch tier {
		case "alpha", "beta", "vip1", "vip2", "vip3":
			return Principal{Role: role, Tier: tier}, nil
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
