package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Benjylem/qcm-forge/internal/config"
	"github.com/Benjylem/qcm-forge/internal/db"
	"github.com/Benjylem/qcm-forge/internal/httpserver"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPool(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	mux := httpserver.NewMux(pool)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
