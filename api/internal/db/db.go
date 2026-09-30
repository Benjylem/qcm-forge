package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool crée un pool de connexions PostgreSQL. La connexion réelle
// n'a pas encore lieu ici : pgxpool.New se contente de valider le DSN
// et se connecte à la demande. C'est Ping() (ailleurs) qui vérifiera
// que la base répond vraiment.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
}
