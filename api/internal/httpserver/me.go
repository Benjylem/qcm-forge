package httpserver

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Benjylem/qcm-forge/internal/auth"
)

// meHandler est la première route "protégée" : elle prouve que le
// cycle complet (cookie -> jeton -> session valide en base) fonctionne.
// Les futures routes privées suivront le même schéma.
func meHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "non connecté")
			return
		}

		userID, err := auth.ValidateSession(r.Context(), pool, cookie.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "non connecté")
			return
		}

		writeJSON(w, http.StatusOK, map[string]int64{"user_id": userID})
	}
}
