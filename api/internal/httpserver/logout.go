package httpserver

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Benjylem/qcm-forge/internal/auth"
)

func logoutHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pas de cookie, ou déjà expiré : le résultat voulu (plus de
		// session active) est déjà atteint, pas la peine de renvoyer
		// une erreur pour ça.
		if cookie, err := r.Cookie(sessionCookieName); err == nil {
			_ = auth.DeleteSession(r.Context(), pool, cookie.Value)
		}
		clearSessionCookie(w)
		w.WriteHeader(http.StatusNoContent)
	}
}
