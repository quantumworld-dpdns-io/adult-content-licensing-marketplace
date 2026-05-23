package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/audit"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/config"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/pkg/httpx"
)

func main() {
	cfg := config.Load()
	_ = os.Setenv("AUTH_SECRET", cfg.AuthSecret)

	repo, err := license.NewRepository(cfg.LicenseStoreBackend, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("repository init failed: %v", err)
	}
	auditStore, err := audit.NewStore(cfg.LicenseStoreBackend, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("audit store init failed: %v", err)
	}

	mux := httpx.NewMux(repo, auditStore)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("gateway listening on :%s (license backend=%s)", cfg.Port, cfg.LicenseStoreBackend)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
