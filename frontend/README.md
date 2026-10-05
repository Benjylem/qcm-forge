# frontend

React + TypeScript, servi par Vite.

```bash
npm ci          # installe les dépendances à partir de package-lock.json
npm run dev     # serveur de dev sur http://localhost:5173
npm run build   # vérification TypeScript + build dans dist/
```

En dev, Vite redirige `/health` et `/api/*` vers l'API Go (`localhost:8080`, voir `vite.config.ts`). Le navigateur ne parle qu'à une seule origine, donc pas besoin de configurer CORS. L'API doit tourner (`docker compose up -d`).
