package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Status != "ok" || got.Service != "gateway" {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestReadyz(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestCreateLicenseForbiddenForLowTier(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	createBody := []byte(`{"id":"lic_1","creator_id":"u_1","title":"Sample","currency":"USDC"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(createBody))
	createReq.Header.Set("X-Role", "regularuser")
	createReq.Header.Set("X-Tier", "alpha")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden, got %d", createRR.Code)
	}
}

func TestCreateLicenseForVIPWithCrypto(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	createBody := []byte(`{"id":"lic_2","creator_id":"u_2","title":"Sample","currency":"USDC"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(createBody))
	createReq.Header.Set("X-Role", "regularuser")
	createReq.Header.Set("X-Tier", "vip1")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected created, got %d", createRR.Code)
	}
}

func TestCreateLicenseRejectsFiat(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	createBody := []byte(`{"id":"lic_3","creator_id":"u_3","title":"Sample","currency":"USD"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(createBody))
	createReq.Header.Set("X-Role", "admin")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", createRR.Code)
	}
}

func TestPolicyEndpointsAuditorAccess(t *testing.T) {
	t.Parallel()
	mux := NewMux()

	req := httptest.NewRequest(http.MethodGet, "/v1/policy/compliance", nil)
	req.Header.Set("X-Role", "auditor")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/v1/policy/threat-model", nil)
	req2.Header.Set("X-Role", "auditor")
	rr2 := httptest.NewRecorder()
	mux.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr2.Code)
	}
}
