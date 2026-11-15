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

## Registration correction and regression review — 15 November 2026

This section records the follow-up QA review that corrected the confirmed registration defect and re-ran the full frontend verification set. Earlier sections are retained unchanged.

### Registration validation defect and correction

- **Area:** Farmer registration (`frontend/src/App.tsx`)
- **Status:** Confirmed defect — **corrected in the frontend**; regression test now passes.
- **Root cause:** The registration form submitted directly to `api.register()` without checking that at least one of email or phone was present. HTML `required` and `minLength` attributes covered the name and password fields in a browser, but the contact-field rule was not enforced by the application.
- **Correction:** Added minimal client-side validation in the registration submit handler:
  1. Full name is required (trimmed value).
  2. Password is required and must contain at least eight characters.
  3. At least one of email or phone is required (trimmed values).
  4. Invalid submissions return before `api.register()` is called.
  5. A clear message is rendered inside the registration dialog with `role="alert"`.
  6. Email-only and phone-only registrations proceed with the expected payload.
- **Design and dependencies:** No interface redesign, no new package, and no backend change.

### Tests added or improved

`frontend/tests/registration.test.tsx` was expanded from one regression test to five:

1. Missing email **and** phone prevents submission and shows a validation message.
2. A provided email allows submission and the payload contains the expected name, email, password and language, with no phone.
3. A provided phone allows submission and the payload contains the expected name, phone, password and language, with no email.
4. A password shorter than eight characters prevents submission and shows a validation message.
5. A missing full name prevents submission and shows a validation message.

The tests use the existing Vitest and React Testing Library setup. Explicit `cleanup()` was added because Vitest globals are not enabled, so automatic cleanup was not registered between cases. Assertions were not weakened and existing tests were not deleted.

### Test execution results

| Check | Command | Result |
|---|---|---|
| Baseline reproduction (before correction) | `npm test` | **FAIL** — 1 failed, 1 passed. `api.register` was called once with both `email` and `phone` undefined. |
| Frontend test suite (after correction) | `npm test` | **PASS** — 2 test files passed, 6 tests passed. |
| Type check | `npm run typecheck` | **PASS** — no TypeScript errors. |
| Production build | `npm run build` | **PASS** — `tsc -b && vite build`; 1580 modules transformed; `dist/assets/index-02JPG0CK.js` 217.92 kB (gzip 68.32 kB). |

### Dependency audit results

| Check | Command | Result |
|---|---|---|
| Full audit | `npm audit` | **FAIL (exit 1)** — 3 vulnerabilities: 1 moderate, 2 critical. All are development/test dependencies. |
| Production-only audit | `npm audit --omit=dev` | **PASS** — 0 vulnerabilities. |

Details:

- `@vitest/mocker <=4.1.11` — **moderate** — advisories GHSA-82fw-gwwq-j7x9 (path traversal / arbitrary file read via redirect mock) and GHSA-qxg6-89mg-rqv8 (code injection via unauthenticated manual mock registration).
- `tinypool <=2.1.1` — **critical** (two advisories: GHSA-5gmw-xhrv-c9v3 and GHSA-85c8-ppgw-ccpr) — prototype-pollution gadget leading to remote code execution in worker options.
- npm reports the fix as `vitest@5.1.0`, a **breaking major upgrade**. `npm audit fix --force` was **not** run. Remediation requires a planned test-tooling upgrade and is outstanding.

### API client contract review (Task D)

`frontend/src/lib/api.ts` was compared against the backend routes and DTOs in this repository. The following were checked and **no endpoint path, request payload, or success-envelope mismatch was found**:

- `POST /farmers/` — matches `CreateFarmerRequest`; response `data` is the farmer object.
- `GET /farms/?user_id=...&page=...&limit=...` — matches `FarmListResult` (`data`, `total`, pagination fields).
- `POST /farms/` — payload matches `CreateFarmRequest`.
- `GET /crops/?page=...&limit=...` — matches the crop list envelope.
- `GET /farms/:id/soil` and `POST /soil/` — match the soil list and create contracts.
- `POST /earthdata/climate` — frontend reads `data.observations` and `data.total_records`, which match `EarthdataResponse`.
- `POST /rotation/` — frontend reads `data.results`, which matches the rotation handler.
- `GET /farms/:id/rotations` — matches the farm-plan list route.
- Error handling reads the top-level `message` field used by the backend error envelope.

No API client change was made, because no mismatch was clearly proven and no live runtime check against a running backend/database was possible in this review. Unresolved integration risks are listed below.

### Authentication limitations (confirmed)

- The frontend has **no login, password-reset, or authenticated-session implementation**.
- `api.ts` **never attaches an `Authorization: Bearer` header**.
- The profile stored in `localStorage` (`soilcycle-profile`) is only a convenience copy of the farmer record, **not** authenticated state and not a security control.
- The backend has JWT validation middleware (`backend/internal/middleware/auth.go`) but no login endpoint or registration flow that issues JWT access tokens.

### Backend route registration inconsistency (confirmed by code inspection)

- The documented entry point `backend/main.go` (`go run .`) registers crops, farmers, farms, soil, rotation and earthdata routes but **does not apply `middleware.RequireAuth()`**. Consequently, the API routes registered through this entry point are not protected by that authentication middleware.
- The alternative `internal/server.New()` applies `RequireAuth()` to the API group but registers **only** crops, farmers and farms — soil, rotation and earthdata would be absent if the server were started via that path.
- **Route registration count:** Inspection of `backend/main.go` confirms that crops, farmers, farms, soil, rotation and earthdata are each registered once. No duplicate module registrations were found in this entry point.

### Farm ownership and soil authorization concerns (code-level, not runtime-exploited in this review)

These were identified by reading the backend service/repository layers. No running database was available, so they are **code-level security concerns requiring backend verification**, not runtime-confirmed exploits:

1. **Farm listing without an owner filter:** `FarmRepository.List` uses `($1::uuid IS NULL OR user_id = $1)`. The handler passes the optional `user_id` query value, so omitting `user_id` returns **all farms**.
2. **Requester-supplied identity for farm update/delete:** both handlers take `user_id` from the query string and pass it to ownership-filtered SQL. The SQL checks the supplied ID, but an unauthenticated caller can supply another farmer's ID; this is not proof of identity.
3. **Farm detail access:** `GET /farms/:id` has no ownership filter.
4. **Soil records:** `GET /soil/:id`, `PUT /soil/:id` and `DELETE /soil/:id` contain no ownership check. Creating a soil record checks only that the farm exists, not that the caller owns it. `GET /farms/:id/soil` returns another farm's records, including exact pH and nutrient measurements, to anyone who knows/guesses the farm ID.
5. **Farmer records:** `GET /farmers/`, `GET /farmers/:id` and `PUT /farmers/:id` have no ownership or role check in the inspected layers; listing exposes names, emails and phone numbers.

These findings must be resolved by the backend developers. QA did not modify backend code, migrations, ownership logic, or trust frontend-supplied IDs.

### Status classification

**Confirmed defects**

- Registration accepted a submission with neither email nor phone. **Corrected in frontend; regression test passes.**
- No usable authentication/session flow exists end-to-end (login endpoint absent, no bearer token sent). **Open.**

**Runtime-verified failures**

- Baseline registration regression test failed before the correction (`api.register` called with both contact fields empty).
- `npm audit` exits 1 with 3 development/test dependency vulnerabilities. Production audit is clean.

**Code-level security concerns (backend)**

- Unauthenticated startup path (`main.go`).
- Route registration inconsistency between `main.go` and `internal/server`.
- Farm listing, farm detail, farm update/delete identity handling.
- Soil read/update/delete/create ownership checks.
- Farmer list/detail/update access.

**Requires further investigation**

- Whether authentication should be enabled via `internal/server` after a login/token-issuance module exists.
- Live frontend/backend contract testing with PostgreSQL and the NASA POWER dependency.
- Browser-level validation behaviour (native constraints plus the new inline messages).
- CORS and deployment configuration against the actual frontend origin.

### Files changed in this review

- `frontend/src/App.tsx` — registration validation and inline validation message.
- `frontend/tests/registration.test.tsx` — expanded regression coverage.
- `docs/TEST_REPORT.md` — this section.

No backend application code, migrations, schemas, dependency manifests or lockfiles were modified. The existing registration test file was expanded with additional regression coverage. No unrelated tests were removed. Nothing was committed or pushed.

### Recommended next QA activities

1. Backend developers to reconcile the two startup paths, ensure all intended routes are registered, and apply authentication consistently; run the backend test suite with a compatible Go toolchain.
2. Add authentication (login/token issuance) before enabling `RequireAuth`; update the frontend API client to send bearer tokens and add session-expiry tests.
3. Add authorization tests proving a farmer cannot read or modify another farmer's farms, soil records, or profile.
4. Run live integration tests against PostgreSQL and the NASA POWER endpoint for registration → farm → soil → climate → rotation flows.
5. Plan the breaking Vitest upgrade that remediates the development dependency advisories, then re-run the full audit and suite.
