# qcm-forge

Application web qui transforme des supports de cours (PDF, PPTX) en fiches de révision et en QCM, avec suivi des erreurs et organisation par matière. Chaque question renvoie à l'extrait du cours d'où elle vient.

Projet personnel, en solo. L'objectif est de montrer un pipeline complet : API, base de données, conteneurs, CI/CD, déploiement et observabilité. Le LLM n'est qu'une brique parmi d'autres.

## État actuel

- [x] PostgreSQL 17 en local via Docker Compose, avec healthcheck
- [x] API Go avec `/health` (build multi-stage, image finale sur Alpine, binaire statique)
- [x] `/health` vérifie réellement la connexion à Postgres (`pool.Ping`) : `200 ok` ou `503 db unreachable`
- [x] Frontend React + Vite minimal qui affiche l'état de `/health` (proxy Vite vers l'API)
- [ ] Fonctionnalités MVP (comptes, upload, extraction, génération de QCM, mode révision)
- [ ] CI GitHub Actions (tests, build, scan Trivy)
- [ ] Déploiement VPS avec HTTPS
- [ ] Terraform + monitoring (Prometheus/Grafana)
- [ ] Test de charge (k6)

## Architecture (état actuel)

```
┌───────────────────────────────────────────────────────┐
│  Docker Compose (réseau interne qcm-forge_default)     │
│                                                         │
│  ┌────────────────────────┐   ┌──────────────────────┐ │
│  │ api (Go, build local)  │   │ db (postgres:17)      │ │
│  │ 127.0.0.1:8080 -> 8080 │──▶│ 127.0.0.1:5433 -> 5432│ │
│  │ GET /health -> "ok"    │   │ volume: pgdata        │ │
│  │ depends_on: db healthy │   │ healthcheck pg_isready│ │
│  └────────────────────────┘   └──────────────────────┘ │
└───────────────────────────────────────────────────────┘
```

- `api` est buildée depuis `api/Dockerfile` en multi-stage : compilation dans `golang:1.26.8-alpine3.24`, puis binaire statique (`CGO_ENABLED=0`) copié dans `alpine:3.24.2` et exécuté en utilisateur `nobody`.
- L'API joint Postgres via le réseau interne Docker (`db:5432`). Le port `5433` côté hôte ne sert qu'à se connecter depuis la machine (psql, DBeaver).
- `depends_on: condition: service_healthy` : l'API ne démarre que quand `pg_isready` répond.
- Tous les ports sont exposés sur `127.0.0.1` uniquement, donc rien n'est accessible depuis le réseau local.

## Prérequis

- Docker Engine et Docker Compose v2 (`docker compose`, sans tiret)
- `curl` pour tester l'API

## Lancer le projet en local

```bash
cp .env.example .env        # puis remplacer le mot de passe
docker compose up -d --build
docker compose ps           # db doit être "healthy", api "running"
curl -i http://127.0.0.1:8080/health
```

Réponse attendue : `HTTP/1.1 200 OK` avec le corps `ok`. Si Postgres est arrêté (`docker compose stop db`), la réponse devient `503 db unreachable`.

### Frontend (dev)

Prérequis supplémentaire : Node.js et npm.

```bash
cd frontend
npm ci
npm run dev                 # http://localhost:5173 affiche "API : ok"
```

Le proxy Vite redirige `/health` et `/api/*` vers `127.0.0.1:8080`, ce qui évite d'avoir à configurer CORS en dev.

Arrêter : `docker compose down`. Pour supprimer aussi les données, ajouter `-v` : le volume `pgdata` est alors détruit.

## Variables d'environnement

Fichier `.env`, jamais commité (voir `.gitignore`). `.env.example` sert de modèle.

| Variable | Utilisée par | Rôle | Défaut |
|---|---|---|---|
| `POSTGRES_USER` | db, api | utilisateur Postgres | — |
| `POSTGRES_PASSWORD` | db, api | mot de passe Postgres | — |
| `POSTGRES_DB` | db, api | nom de la base | — |
| `POSTGRES_HOST` | api | hôte de la base | `db` (nom du service Compose) |
| `POSTGRES_PORT` | api | port de la base | `5432` (port interne au conteneur) |

## Structure du dépôt

```
api/
  cmd/api/main.go          point d'entrée : config -> pool pgx -> serveur HTTP
  internal/config/         lecture des variables d'environnement
  internal/db/             création du pool de connexions pgx
  internal/httpserver/     routes HTTP (/health)
  internal/auth/           (à venir) comptes utilisateurs
  internal/upload/         (à venir) upload et extraction PDF/PPTX
  internal/llmprovider/    (à venir) interface vers le fournisseur LLM
  internal/qcm/            (à venir) génération et stockage des QCM
  Dockerfile
frontend/                  React + TypeScript (Vite), affiche l'état de /health
migrations/                (à venir) migrations SQL
docker-compose.yml
```

## Décisions techniques

| Choix | Pourquoi |
|---|---|
| **Go** pour l'API | Binaire statique, donc image Docker légère. Concurrence native pour le traitement des documents. Docker, Kubernetes et Terraform sont écrits en Go. |
| **PostgreSQL** plutôt que SQLite | L'API et un futur worker écriront en parallèle depuis des conteneurs séparés. JSONB pour les QCM et recherche plein texte pour les extraits. SQLite suffirait pour 3 utilisateurs, mais Postgres permet d'apprendre l'exploitation d'une base client-serveur. |
| **Clé API LLM fournie par chaque utilisateur** | Chaque utilisateur paie sa propre consommation. Les clés seront chiffrées en base (AES-GCM), jamais renvoyées au frontend ni écrites dans les logs. |
| **Fournisseur LLM derrière une interface Go** | Pas de dépendance à un seul fournisseur. |
| **Versions d'images fixées** | Builds reproductibles. Pas de `latest`. |

## Notes Fedora / SELinux

Pour monter un **dossier local** dans un conteneur, ajouter `:z` au volume (ex. `./migrations:/migrations:z`), sinon SELinux bloque l'accès. Ce n'est pas nécessaire pour les volumes nommés comme `pgdata`.

## Stack

- **API** : Go, `net/http`, `pgx`
- **Frontend** : React, TypeScript, Vite
- **Base de données** : PostgreSQL 17
- **Conteneurs** : Docker, Docker Compose
