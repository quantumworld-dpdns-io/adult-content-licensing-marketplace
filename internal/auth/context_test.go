package auth

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestPrincipalFromRequestRegularUserVIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Role", "regularuser")
	r.Header.Set("X-Tier", "vip2")
	p, err := PrincipalFromRequest(r)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.Role != "regularuser" || p.Tier != "vip2" {
		t.Fatalf("unexpected principal: %+v", p)
	}
}

func TestPrincipalFromBearer(t *testing.T) {
	t.Setenv("AUTH_SECRET", "abc123")
	tok, err := SignToken("abc123", Claims{Sub: "u1", Role: "regularuser", Tier: "vip1", Exp: time.Now().Add(time.Hour).Unix()})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+tok)
	p, err := PrincipalFromRequest(r)
	if err != nil {
		t.Fatalf("principal: %v", err)
	}
	if p.Role != "regularuser" || p.Tier != "vip1" || p.Sub != "u1" {
		t.Fatalf("unexpected principal: %+v", p)
	}
}

func TestCanCreateLicense(t *testing.T) {
	_ = os.Setenv("AUTH_SECRET", "test")
	if !CanCreateLicense(Principal{Role: "admin"}) {
		t.Fatal("admin should be allowed")
	}
	if CanCreateLicense(Principal{Role: "regularuser", Tier: "alpha"}) {
		t.Fatal("alpha should not create")
	}
	if !CanCreateLicense(Principal{Role: "regularuser", Tier: "vip1"}) {
		t.Fatal("vip1 should create")
	}
}
