package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// healthHandler répond "ok" si la base répond, ou 503 sinon. Le but
// est qu'un outil de supervision (plus tard) puisse distinguer
// "l'API tourne mais la base est injoignable" d'un vrai "tout va bien".
func healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Timeout court : /health doit répondre vite. Si la base met
		// plus de 2s à répondre à un ping, ce n'est pas sain non plus.
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("db unreachable"))
			return
		}

		w.Write([]byte("ok"))
	}
}
