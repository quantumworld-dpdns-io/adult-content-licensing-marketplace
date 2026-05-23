package auth

import (
	"net/http/httptest"
	"testing"
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

func TestCanCreateLicense(t *testing.T) {
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
