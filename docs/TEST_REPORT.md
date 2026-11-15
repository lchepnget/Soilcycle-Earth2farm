# Verification report

Date: 2026-11-14

| Check | Command / method | Result | Notes |
|---|---|---|---|
| Archive inspection | Python `zipfile` listing and extraction | Passed | 118 archive entries including Git metadata; 55 non-Git project files extracted from the supplied archive. |
| Original file preservation | SHA-256 inventory comparison | Passed | 53 of 54 original non-Git files are byte-identical; `.gitignore` is the sole modified original file; no original files were removed. The inventory JSON is included in `docs/reports/`. |
| Backend tests (baseline) | `cd backend && GOTOOLCHAIN=local go test ./...` | Blocked | Local Go is 1.23.2; `go.mod` requires Go 1.26.8. |
| Frontend dependency installation | `cd frontend && npm install --fetch-retries=0 --fetch-timeout=20000` | Blocked | npm registry DNS lookup failed (`EAI_AGAIN registry.npmjs.org`). |
| Frontend type check | `tsc --noEmit --pretty false` | Blocked | Dependencies and React/Vite type declarations are unavailable because installation could not complete. The compiler output is therefore dominated by missing module/type declarations; this is not a successful type check. |
| Frontend unit tests | `npm test` | Not run | Vitest and Testing Library could not be installed. |
| Frontend production build | `npm run build` | Not run | Vite/TypeScript dependencies could not be installed. |
| Backend/frontend live integration | Requires PostgreSQL, supported Go toolchain and running API | Not run | No database service or compatible Go toolchain was available in the execution environment. |
| ZIP integrity | Python ZIP CRC test and clean temporary extraction | Passed | Archive contains 70 entries; CRC test passed; complete root directory extracted and backend, migration, frontend, tests and docs were verified present. |

## Important release blockers

- A `package-lock.json` could not be generated without npm registry access. The package manifest pins direct dependency versions; run `npm install` in a network-enabled environment and commit the resulting lockfile before treating installs as reproducible.
- Backend tests and runtime verification require Go 1.26.8 or the toolchain version declared by the project. The supplied environment has Go 1.23.2.
- Authentication and ownership enforcement gaps are documented in `INTEGRATION_AUDIT.md`. Do not deploy this project with real farmer data until those backend gaps are resolved.
