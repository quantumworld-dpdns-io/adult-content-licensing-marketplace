package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if !driverRegistered("postgres") {
		log.Fatal("postgres driver is not registered in this binary")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := license.OpenPostgresAndEnsureSchema(ctx, dsn)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	defer db.Close()

	fmt.Println("migration complete")
}

func driverRegistered(name string) bool {
	for _, d := range sql.Drivers() {
		if d == name {
			return true
		}
	}
	return false
}
