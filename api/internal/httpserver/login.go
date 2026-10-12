package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Benjylem/qcm-forge/internal/auth"
)

func loginHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			writeError(w, http.StatusBadRequest, "corps de requête invalide")
			return
		}

		user, err := auth.GetUserByEmail(r.Context(), pool, creds.Email)
		if errors.Is(err, auth.ErrUserNotFound) {
			// Même message que "mauvais mot de passe" : ne jamais
			// révéler si un email existe déjà, ça faciliterait
			// l'énumération de comptes par un attaquant.
			writeError(w, http.StatusUnauthorized, "email ou mot de passe incorrect")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		ok, err := auth.CheckPassword(creds.Password, user.PasswordHash)
		if err != nil || !ok {
			writeError(w, http.StatusUnauthorized, "email ou mot de passe incorrect")
			return
		}

		token, err := auth.CreateSession(r.Context(), pool, user.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		setSessionCookie(w, token)
		w.WriteHeader(http.StatusOK)
	}
}
