package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Benjylem/qcm-forge/internal/auth"
)

// credentials est partagé par /signup et /login : même forme de
// requête dans les deux cas.
type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func signupHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			writeError(w, http.StatusBadRequest, "corps de requête invalide")
			return
		}

		// Validation minimale : un vrai contrôle d'email (RFC 5322) est
		// un piège classique, inutile ici. "@" suffit à filtrer les
		// erreurs de saisie grossières.
		if !strings.Contains(creds.Email, "@") {
			writeError(w, http.StatusBadRequest, "email invalide")
			return
		}
		if len(creds.Password) < 8 {
			writeError(w, http.StatusBadRequest, "mot de passe trop court (8 caractères minimum)")
			return
		}

		hash, err := auth.HashPassword(creds.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		userID, err := auth.CreateUser(r.Context(), pool, creds.Email, hash)
		if errors.Is(err, auth.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email déjà utilisé")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		token, err := auth.CreateSession(r.Context(), pool, userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		setSessionCookie(w, token)
		w.WriteHeader(http.StatusCreated)
	}
}
