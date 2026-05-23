package payment

import "testing"

func TestEnsureCryptoCurrency(t *testing.T) {
	if err := EnsureCryptoCurrency("USDC"); err != nil {
		t.Fatalf("expected USDC allowed: %v", err)
	}
	if err := EnsureCryptoCurrency("USD"); err == nil {
		t.Fatal("expected USD rejected")
	}
}
