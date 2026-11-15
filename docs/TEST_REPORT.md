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

## Follow-up verification — 15 November 2026

### Frontend verification results

| Check | Command | Result | Notes |
|---|---|---|---|
| Dependency installation | `npm ci` | PASS | Dependencies installed successfully. |
| Baseline unit tests | `npm test` | PASS | Original suite: 1 test passed. |
| Type checking | `npm run typecheck` | PASS | No TypeScript errors. |
| Production build | `npm run build` | PASS | Production build completed. |
| Production dependency audit | `npm audit --omit=dev` | PASS | Zero reported vulnerabilities. |
| Full dependency audit | `npm audit` | FAIL | Three vulnerabilities in development/test dependencies: one moderate and two critical. |
| Registration regression test | `npx vitest run tests/registration.test.tsx` | FAIL — defect reproduced | Registration API was called despite both email and phone being empty. |

### Confirmed defect: Missing registration validation

- **Area:** Farmer registration
- **Suggested severity:** Medium
- **Status:** Reproduced; awaiting correction.

The backend contract in `backend/internal/modules/farmers/dto.go` requires either an email address or a phone number.

The frontend currently permits both fields to be empty and still submits the registration request.

The regression test in `frontend/tests/registration.test.tsx` reproduced the issue: `api.register()` was called once when it should not have been called.

**Expected behaviour:** The frontend should prevent submission and display a clear validation message when both contact fields are empty.

**Recommended action:** Frontend developers should implement this validation. QA will rerun the regression test after the correction.

### Security findings

The production dependency audit reported zero vulnerabilities. The full audit reported three vulnerabilities in development/test dependencies, including two critical findings. Developer review and remediation are outstanding. QA has not applied forced dependency upgrades.

### Outstanding QA actions

1. Frontend developers to correct registration validation.
2. QA to rerun the regression test.
3. QA to rerun the full frontend test suite, type check and production build.
4. Developers to review the dependency audit findings.
5. Complete live frontend/backend integration testing.

### QA conclusion

The frontend baseline tests, type check and production build passed. However, targeted testing identified a reproducible registration validation defect. Development dependency vulnerabilities also remain outstanding.

**Recommendation:** The registration workflow should not be considered fully verified until the defect is corrected and the regression test passes.
