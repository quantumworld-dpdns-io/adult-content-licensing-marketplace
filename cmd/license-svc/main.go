package main

import (
	"log"
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/pkg/api"
)

func main() {
	store := license.NewStore()
	h := api.LicenseHandler{Store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/licenses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.List(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	log.Println("license-svc listening on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
