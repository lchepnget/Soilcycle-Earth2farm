# SoilCycle repository audit and frontend integration

## Repository as supplied

The supplied archive contains a Go backend under `backend/`, PostgreSQL schema/migration and crop seed SQL under `backend/internal/database/`, project documentation under `docs/`, and a root README. It does not contain a frontend directory, frontend package manifest, frontend tests, Docker configuration, CI workflow, or root database directory. The archive also contains Git metadata; it is excluded from the delivery archive to avoid shipping repository internals.

## Backend architecture observed

- Language/framework: Go with Fiber v3.
- Database: PostgreSQL accessed through `pgxpool`.
- Configuration: environment variables loaded by `godotenv` (`APP_PORT`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`).
- Entrypoint: `backend/main.go`.
- Schema: one SQL migration creates users, farms, crops, crop characteristics, farm crop history, soil data, climate data, rotation plans and rotation results.
- Password registration: farmer creation hashes the supplied password with bcrypt; no login/token-issuing endpoint is implemented.
- Authentication middleware: JWT HS256 validation exists in `backend/internal/middleware/auth.go`; the middleware fails closed if `JWT_SECRET` is missing or shorter than 32 characters. The existing `server.New` / `RegisterRoutes` path applies this middleware to `/api/v1` routes, while the original `main.go` constructs a separate Fiber app and does not apply it.
- External data: the Earthdata module calls NASA POWER daily point API for precipitation, temperature and humidity parameters. The code does not demonstrate integrations for GPM, SMAP, MODIS or SoilGrids.

## Actual registered API routes

The frontend API client is based on route registration and handlers, not the proposal examples in `docs/API_CONTRACTS.md`.

| Method | Actual route | Current use |
|---|---|---|
| POST | `/api/v1/farmers/` | Farmer registration |
| GET | `/api/v1/farmers/` | Farmer list |
| GET | `/api/v1/farmers/:id` | Farmer lookup |
| PUT | `/api/v1/farmers/:id` | Farmer update |
| GET/POST/PUT/DELETE | `/api/v1/farms/` and `/api/v1/farms/:id` | Farm CRUD |
| GET/POST/PUT/DELETE | `/api/v1/crops/` and `/api/v1/crops/:id` | Crop catalogue and CRUD |
| GET/PUT | `/api/v1/crops/:id/characteristics` | Crop characteristics |
| GET/POST/PUT/DELETE | `/api/v1/soil/` and `/api/v1/soil/:id` | Soil record CRUD |
| GET | `/api/v1/farms/:id/soil` | Soil records by farm |
| POST | `/api/v1/rotation/` | Generate and save rotation plan |
| GET | `/api/v1/rotation/:id` | Retrieve rotation plan |
| GET | `/api/v1/rotation/:id/results` | Retrieve rotation results |
| GET | `/api/v1/farms/:id/rotations` | List farm rotation plans |
| POST | `/api/v1/earthdata/climate` | Retrieve NASA POWER climate observations |

The API contract document's `/farms/:id/crops`, `/farms/:id/soil-data`, `/farms/:id/environment`, `/recommendations`, and `/scenarios/compare` examples are proposals and are not registered routes in the supplied source.

## Frontend implementation

A React + TypeScript + Vite frontend has been added under `frontend/`. It provides a public onboarding state and views for overview, farm profiles, soil records, climate observations, crop rotation, crop library, and help. The API calls use the actual route paths and request fields found in the Go DTOs and handlers. Empty states are used instead of synthetic farm counts, soil measurements, climate series, or crop results.

## Known backend integration and security gaps

1. The backend has farmer registration but no login, logout, password reset, or JWT issuance endpoint. The UI therefore does not claim to provide authenticated sessions. The browser stores only the returned profile record for convenience; this is not authentication.
2. The original `backend/main.go` creates a Fiber app without attaching the existing JWT middleware. The alternate `internal/server.New` path applies JWT middleware, but it is not the original executable entrypoint. This discrepancy must be resolved before production deployment.
3. Farm create accepts `user_id` in the request body; farm update/delete accept `user_id` as a query parameter. This is explicitly marked temporary in the backend and is not a secure ownership boundary.
4. Several soil, crop, rotation and farm operations need consistent ownership/role checks before multi-user production use.
5. Crop rotation logic exists, but its agronomic assumptions require domain expert review; UI copy treats its results as decision support rather than validated predictions.
6. NASA POWER is integrated in code. Other datasets described as potential sources in project documents are not confirmed integrations in the supplied source.
7. No dedicated soil-health scoring endpoint or laboratory integration exists. The UI records soil values and source labels without calculating a health score.
8. No farm history/crop planting endpoint is registered in the supplied routes, so the frontend does not offer farm crop-history entry.

## Preservation

The implementation is additive: the original backend, SQL migration/seed, documentation and source files are retained. Root `.gitignore` was extended for frontend dependency/build artifacts. New frontend files and this audit documentation were added. No original source files were removed or modified.

The original `backend/main.go` also registers the crops module three times and farmers/farms more than once. This was left unchanged to preserve the supplied backend; it should be cleaned up and covered by route-registration tests in a backend maintenance change. The `/health` handler in that entrypoint returns a healthy response without performing a database ping, unlike the alternate `internal/server.RegisterRoutes` health handler.
