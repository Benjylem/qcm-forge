package httpserver

import (
	"net/http"
	"time"

	"github.com/Benjylem/qcm-forge/internal/auth"
)

const sessionCookieName = "session"

// setSessionCookie pose le cookie de session.
//
// Secure: false pour l'instant, en connaissance de cause : on tourne
// en HTTP en local (pas de certificat TLS, ça viendra à l'étape 4 -
// déploiement VPS). Un cookie de session DOIT être Secure en
// production, sinon il circule en clair sur le réseau (wifi public
// compris) et peut être volé. À changer en true avant tout déploiement
// réel - ne pas oublier.
func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(auth.SessionTTL),
	})
}

// clearSessionCookie écrase le cookie avec MaxAge négatif, ce qui dit
// au navigateur de le supprimer immédiatement.
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
