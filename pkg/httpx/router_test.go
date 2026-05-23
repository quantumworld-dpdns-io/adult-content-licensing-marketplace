package httpx

import (
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/audit"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	t.Parallel()
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())
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
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestCreateLicenseForbiddenForLowTier(t *testing.T) {
	t.Parallel()
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())
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
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())
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
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())
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
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())

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

func TestDevTokenAndBearerAccess(t *testing.T) {
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())

	t.Setenv("AUTH_SECRET", "local-secret")
	mintReq := httptest.NewRequest(http.MethodPost, "/v1/auth/dev-token?role=regularuser&tier=vip1&sub=user-9", nil)
	mintRR := httptest.NewRecorder()
	mux.ServeHTTP(mintRR, mintReq)
	if mintRR.Code != http.StatusOK {
		t.Fatalf("expected token mint 200, got %d", mintRR.Code)
	}

	var tokenResp map[string]string
	if err := json.Unmarshal(mintRR.Body.Bytes(), &tokenResp); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	token := tokenResp["token"]
	if token == "" {
		t.Fatal("expected token")
	}

	createBody := []byte(`{"id":"lic_tok_1","creator_id":"u_9","title":"Token Flow","currency":"ETH"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected created with token, got %d", createRR.Code)
	}
}

func TestAuditEventsEndpoint(t *testing.T) {
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())

	createBody := []byte(`{"id":"lic_aud_1","creator_id":"u_20","title":"Audit","currency":"USDC"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(createBody))
	createReq.Header.Set("X-Role", "admin")
	createReq.Header.Set("X-Sub", "admin-1")
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected created, got %d", createRR.Code)
	}

	auditReq := httptest.NewRequest(http.MethodGet, "/v1/audit/events", nil)
	auditReq.Header.Set("X-Role", "auditor")
	auditRR := httptest.NewRecorder()
	mux.ServeHTTP(auditRR, auditReq)
	if auditRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", auditRR.Code)
	}
}

func TestRequestIDHeaderIsSet(t *testing.T) {
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header")
	}
}

func TestMetricsEndpointAuditorAccess(t *testing.T) {
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())

	seed := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	seedRR := httptest.NewRecorder()
	mux.ServeHTTP(seedRR, seed)

	req := httptest.NewRequest(http.MethodGet, "/v1/metrics/latency", nil)
	req.Header.Set("X-Role", "auditor")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestAuditEventsEndpointWithFilters(t *testing.T) {
	mux := NewMux(license.NewMemoryRepository(), audit.NewMemoryStore())

	mk := func(id, cur string) {
		body := []byte(`{"id":"` + id + `","creator_id":"u_x","title":"T","currency":"` + cur + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(body))
		req.Header.Set("X-Role", "admin")
		req.Header.Set("X-Sub", "admin-1")
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("seed create failed: %d", rr.Code)
		}
	}
	mk("lic_f1", "USDC")
	mk("lic_f2", "ETH")

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events?type=license.created&entity_id=lic_f2&limit=1&offset=0", nil)
	req.Header.Set("X-Role", "auditor")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	items, ok := out["items"].([]any)
	if !ok {
		t.Fatalf("expected items array, got %T", out["items"])
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 filtered item, got %d", len(items))
	}
	if int(out["total"].(float64)) < 1 {
		t.Fatalf("expected total >= 1, got %v", out["total"])
	}
	if int(out["limit"].(float64)) != 1 {
		t.Fatalf("expected limit 1, got %v", out["limit"])
	}
	if int(out["offset"].(float64)) != 0 {
		t.Fatalf("expected offset 0, got %v", out["offset"])
	}
}
