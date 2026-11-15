# Frontend development and integration

## Requirements

- Node.js 20.19+ or 22.12+ recommended.
- npm 10+.
- The existing Go API and PostgreSQL database for live data operations.

## Install and configure

```bash
cd frontend
cp .env.example .env
npm install
```

Set `VITE_API_BASE_URL` in `frontend/.env` to the reachable backend API base URL. Values prefixed with `VITE_` are public browser configuration; never place API keys, passwords, JWT signing secrets or database credentials there.

## Run the frontend

```bash
npm run dev
```

Vite prints the local development URL. The default API URL is `http://localhost:8090/api/v1`.

## Checks

```bash
npm run typecheck
npm test
npm run build
```

## Existing backend

The backend's `.env.example` is under `backend/`. Configure PostgreSQL using the existing variables, apply `backend/internal/database/migrations/001_initial_schema.sql`, and seed crops using `backend/internal/database/seeds/crops.sql` if desired. Then run the Go backend from `backend/` with the Go version required by `backend/go.mod`.

The current environment used for this delivery had Go 1.23.2 while `backend/go.mod` requires Go 1.26.8. The backend test suite could not be executed in that environment. npm registry access was unavailable during packaging, so dependency installation, generation of `package-lock.json`, frontend type checking, frontend tests and production build could not be completed or represented as passing. Run `npm install` in a network-enabled development environment; commit the resulting lockfile before production use.

## Integration notes

The frontend currently uses the implemented routes documented in `INTEGRATION_AUDIT.md`. The supplied backend lacks a login/token-issuance endpoint, and the main executable does not currently use the repository's JWT middleware path. The frontend onboarding flow is therefore registration/profile setup only, not secure sign-in. Resolve authentication and resource ownership in the backend before exposing this application to real users or production data.
