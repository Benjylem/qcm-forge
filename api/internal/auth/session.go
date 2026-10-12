package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SessionTTL est exporté car httpserver en a besoin pour poser
// l'expiration du cookie, alignée sur celle stockée en base.
const SessionTTL = 7 * 24 * time.Hour

var ErrSessionNotFound = errors.New("auth: session introuvable ou expirée")

// newSessionToken génère un jeton aléatoire opaque de 32 octets, encodé
// en base64 pour pouvoir voyager dans un cookie. crypto/rand (pas
// math/rand) : la source doit être imprédictible, pas juste "varié".
func newSessionToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// hashToken calcule l'empreinte stockée en base. SHA-256 suffit (voir
// pourquoi dans la migration 000002) : le jeton est déjà aléatoire sur
// 32 octets, pas besoin d'un hash lent comme pour un mot de passe.
func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

// CreateSession génère un jeton, l'enregistre (hashé) pour userID, et
// renvoie le jeton en clair : c'est la seule fois où il existe quelque
// part — à mettre dans le cookie immédiatement, il ne sera plus jamais
// récupérable depuis la base.
func CreateSession(ctx context.Context, pool *pgxpool.Pool, userID int64) (string, error) {
	token, err := newSessionToken()
	if err != nil {
		return "", err
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, $3)`,
		hashToken(token), userID, time.Now().Add(SessionTTL),
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateSession renvoie l'userID associé à un jeton encore valide.
func ValidateSession(ctx context.Context, pool *pgxpool.Pool, token string) (int64, error) {
	var userID int64
	err := pool.QueryRow(ctx,
		`SELECT user_id FROM sessions WHERE token_hash = $1 AND expires_at > now()`,
		hashToken(token),
	).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrSessionNotFound
	}
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// DeleteSession supprime une session (déconnexion). Supprimer un jeton
// qui n'existe déjà plus n'est pas une erreur : le résultat recherché
// (plus de session valide) est déjà atteint.
func DeleteSession(ctx context.Context, pool *pgxpool.Pool, token string) error {
	_, err := pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, hashToken(token))
	return err
}
