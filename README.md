# qcm-forge

Application web qui transforme des supports de cours (PDF, PPTX) en fiches de révision et en QCM, avec suivi des erreurs et organisation par matière. Chaque question renvoie à l'extrait du cours d'où elle vient.

## État actuel

- [x] PostgreSQL 17 en local via Docker Compose, avec healthcheck
- [x] API Go avec `/health` (build multi-stage, image finale sur Alpine, binaire statique)
- [ ] Fonctionnalités MVP (upload, extraction, génération de QCM)
- [ ] CI GitHub Actions (tests, build, scan Trivy)
- [ ] Déploiement VPS avec HTTPS
- [ ] Terraform + monitoring

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

`api` est buildée depuis `api/Dockerfile` (multi-stage : compilation dans `golang:alpine`,
binaire statique copié dans une image `alpine` finale). Elle ne contacte pas encore la base
(`/health` vérifie seulement que le processus répond).

## Lancer le projet en local

1. Copier `.env.example` en `.env` et renseigner les variables.
2. `docker compose up -d`
3. Vérifier que le service est en bonne santé : `docker compose ps` (statut `healthy`).

## Stack

- **API** : Go
- **Base de données** : PostgreSQL 17
- **Orchestration locale** : Docker Compose
