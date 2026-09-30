package httpserver

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewMux construit le routeur HTTP de l'API. Le pool de connexions
// est injecté ici plutôt que créé dans le handler, pour rester
// testable et pour ne créer qu'un seul pool partagé par toutes les
// requêtes.
func NewMux(pool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(pool))
	return mux
}
