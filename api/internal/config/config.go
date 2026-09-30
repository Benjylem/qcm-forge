package config

import (
	"fmt"
	"os"
)

// Config regroupe la configuration lue depuis l'environnement.
type Config struct {
	DatabaseDSN string
}

// Load construit la config à partir des variables d'environnement.
// POSTGRES_HOST et POSTGRES_PORT ont des valeurs par défaut adaptées
// au réseau interne de Docker Compose (le service s'appelle "db" et
// expose le port 5432 en interne, indépendamment du port 5433 mappé
// sur l'hôte).
func Load() Config {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := envOrDefault("POSTGRES_HOST", "db")
	port := envOrDefault("POSTGRES_PORT", "5432")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)

	return Config{DatabaseDSN: dsn}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
