# qcm-forge

Application web qui transforme des supports de cours (PDF, PPTX) en fiches de révision et en QCM, avec suivi des erreurs et organisation par matière. Chaque question renvoie à l'extrait du cours d'où elle vient.

## État actuel

- [x] PostgreSQL 17 en local via Docker Compose, avec healthcheck
- [ ] API Go avec `/health`
- [ ] Fonctionnalités MVP (upload, extraction, génération de QCM)
- [ ] CI GitHub Actions (tests, build, scan Trivy)
- [ ] Déploiement VPS avec HTTPS
- [ ] Terraform + monitoring

## Architecture (état actuel)

```
┌─────────────────────────────┐
│  Docker Compose             │
│                             │
│  ┌───────────────────────┐  │
│  │ db (postgres:17)      │  │
│  │ 127.0.0.1:5433 -> 5432│  │
│  │ volume: pgdata        │  │
│  └───────────────────────┘  │
└─────────────────────────────┘
```

## Lancer le projet en local

1. Copier `.env.example` en `.env` et renseigner les variables.
2. `docker compose up -d`
3. Vérifier que le service est en bonne santé : `docker compose ps` (statut `healthy`).

## Stack

- **API** : Go
- **Base de données** : PostgreSQL 17
- **Orchestration locale** : Docker Compose
