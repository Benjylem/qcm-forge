package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

var (
	ErrEmailTaken   = errors.New("auth: email déjà utilisé")
	ErrUserNotFound = errors.New("auth: utilisateur introuvable")
)

// Code Postgres renvoyé en cas de violation d'une contrainte UNIQUE
// (ici, users_email_lower_idx). Évite de faire un SELECT avant chaque
// INSERT pour vérifier l'unicité : la base le garantit déjà, on se
// contente de traduire son erreur.
const pgUniqueViolation = "23505"

func CreateUser(ctx context.Context, pool *pgxpool.Pool, email, passwordHash string) (int64, error) {
	var id int64
	err := pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		email, passwordHash,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return 0, ErrEmailTaken
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (User, error) {
	var u User
	err := pool.QueryRow(ctx,
		`SELECT id, email, password_hash FROM users WHERE lower(email) = lower($1)`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return u, nil
}
