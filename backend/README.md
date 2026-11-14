# SoilCycle by Earth2Farm – Backend

## 1. Project Overview

SoilCycle is a crop-rotation decision-support system developed under Earth2Farm. It aims to help farmers make informed crop-rotation decisions using farm information, soil data, crop characteristics, and environmental data.

The backend provides REST APIs for:

- Farmer management
- Farm management
- Crop management
- Soil data management
- NASA environmental data retrieval
- Crop rotation recommendations

### Technology Stack

| Technology | Purpose |
|---|---|
| Go | Backend programming language |
| Fiber v3 | HTTP framework |
| PostgreSQL 16 | Relational database |
| pgx/v5 | PostgreSQL driver |
| Docker | Database containerization |
| NASA POWER API | Climate and environmental data |

---

## 2. Backend Architecture

```text
                 SOILCYCLE DATABASE
                         |
          +--------------+--------------+
          |              |              |
        USERS           FARMS          CROPS
                         |              |
                     FARM_CROPS         |
                         |              |
              +----------+----------+   |
              |                     |   |
          SOIL_DATA             CLIMATE_DATA
              |                     |
              +----------+----------+
                         |
                   ROTATION_PLANS
                         |
                   ROTATION_RESULTS
```

### Backend Data Flow

```text
Farmer Input
     |
     v
Frontend Application
     |
     v
Go Backend (Fiber)
     |
     +----------------------+
     |                      |
     v                      v
NASA POWER API         PostgreSQL
     |                      |
     +----------+-----------+
                |
                v
       Crop Rotation Engine
                |
                v
       Rotation Recommendations
                |
                v
          Farmer Dashboard
```

---

## 3. Backend Folder Structure

```text
backend/
│
├── main.go
├── go.mod
├── go.sum
├── README.md
│
├── internal/
│   ├── config/
│   ├── database/
│   └── modules/
│       ├── farmers/
│       ├── farms/
│       ├── crops/
│       ├── soil/
│       ├── earthdata/
│       └── rotation/
│
├── migrations/
│
└── .env.example
```

Each module follows a layered structure where applicable:

```text
module/
├── handler.go
├── service.go
├── repository.go
├── model.go
├── dto.go
└── routes.go
```

---

# 4. Starting PostgreSQL with Docker

The development database runs in a Docker container named:

```text
soilcycle-db
```

The container uses PostgreSQL 16.

### Step 1: Start Docker Desktop

Ensure Docker Desktop or the Docker service is running.

### Step 2: Start the PostgreSQL container

```bash
docker start soilcycle-db
```

Expected output:

```text
soilcycle-db
```

### Step 3: Verify the container

```bash
docker ps
```

Confirm that `soilcycle-db` is running.

The development database uses:

| Setting | Value |
|---|---|
| Container | soilcycle-db |
| Image | postgres:16 |
| Host port | 5433 |
| Container port | 5432 |
| Database | soilcycle |
| Database user | soilcycle |

Database credentials should be configured through environment variables and must not be committed to Git.

### Step 4: Check PostgreSQL logs

```bash
docker logs soilcycle-db --tail 30
```

This helps confirm that PostgreSQL started successfully.

---

# 5. Connecting to PostgreSQL

Connect to the database using:

```bash
docker exec -it soilcycle-db psql -U soilcycle -d soilcycle
```

A successful connection displays:

```text
soilcycle=#
```

### Check available tables

Inside PostgreSQL, run:

```sql
\dt
```

The following tables were verified:

```text
climate_data
crop_characteristics
crops
farm_crops
farms
rotation_plans
rotation_results
soil_data
users
```

### Inspect crop tables

```sql
\d crops
```

```sql
\d crop_characteristics
```

### Check crop records

```sql
SELECT id, name, crop_family
FROM crops;
```

### Check crop characteristics

```sql
SELECT
    c.name,
    cc.water_requirement,
    cc.growth_duration_days,
    cc.nitrogen_fixing,
    cc.nutrient_demand
FROM crops c
JOIN crop_characteristics cc
ON c.id = cc.crop_id;
```

Seven crop-characteristics records were inserted and verified during development.

### Exit PostgreSQL

```sql
\q
```

---

# 6. Running the Backend

From the backend directory:

### Download dependencies

```bash
go mod tidy
```

### Format Go files

```bash
gofmt -w .
```

### Run the backend

```bash
go run .
```

The development server runs on:

```text
http://localhost:8090
```

API base URL:

```text
http://localhost:8090/api/v1
```

Ensure your PostgreSQL environment configuration is correct before starting the server.

---

# 7. Backend Validation Tests

The following commands were executed during development.

### Format code

```bash
go fmt ./...
```

### Run Go vet

```bash
go vet ./...
```

### Build the backend

```bash
go build ./...
```

### Run Go tests

```bash
go test ./...
```

### Results

The following checks completed successfully:

- Go formatting
- `go vet`
- Go build
- Go test command execution

**Note:** The Go packages reported `[no test files]`. This means the command completed, but automated Go unit tests have not yet been implemented.

API functionality was also checked manually using `curl`.

---

# 8. Soil Module API Tests

The Soil module supports creating, listing, retrieving, updating, and deleting soil records.

### Available endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/api/v1/soil` | Create soil record |
| GET | `/api/v1/soil` | List soil records |
| GET | `/api/v1/soil/:id` | Retrieve soil record |
| PUT | `/api/v1/soil/:id` | Update soil record |
| DELETE | `/api/v1/soil/:id` | Delete soil record |
| GET | `/api/v1/farms/:id/soil` | List soil records for a farm |

### Test 1: Create soil record

```bash
curl -X POST http://localhost:8090/api/v1/soil \
-H "Content-Type: application/json" \
-d '{
  "farm_id": "192f23b1-aedb-4e78-8af7-81c735fb695c",
  "ph": 6.5,
  "organic_carbon": 1.25,
  "nitrogen": 0.18,
  "phosphorus": 12.5,
  "potassium": 0.35,
  "soil_texture": "Loamy",
  "data_source": "Farmer Entry",
  "measurement_date": "2026-11-14"
}'
```

Expected result:

```json
{
  "message": "Soil record created successfully",
  "success": true
}
```

HTTP status: `201 Created`

### Test 2: List soil records

```bash
curl -s http://localhost:8090/api/v1/soil | python3 -m json.tool
```

Verified that the endpoint returned:

- Soil record data
- Total records
- Current page
- Page limit
- Total pages

HTTP status: `200 OK`

### Test 3: Retrieve soil record by ID

```bash
curl -s http://localhost:8090/api/v1/soil/0ea0aacf-0655-43ba-9214-aa61e7470b09 | python3 -m json.tool
```

HTTP status: `200 OK`

### Test 4: Retrieve soil records by farm

```bash
curl -s http://localhost:8090/api/v1/farms/192f23b1-aedb-4e78-8af7-81c735fb695c/soil | python3 -m json.tool
```

HTTP status: `200 OK`

### Test 5: Update soil record

```bash
curl -X PUT http://localhost:8090/api/v1/soil/0ea0aacf-0655-43ba-9214-aa61e7470b09 \
-H "Content-Type: application/json" \
-d '{
  "ph": 6.8,
  "organic_carbon": 1.35,
  "nitrogen": 0.20,
  "phosphorus": 13.0,
  "potassium": 0.40,
  "soil_texture": "Clay Loam",
  "data_source": "Farmer Entry",
  "measurement_date": "2026-11-14"
}'
```

Expected result:

```json
{
  "message": "Soil record updated successfully",
  "success": true
}
```

HTTP status: `200 OK`

A subsequent GET request confirmed the updated values.

### Test 6: Invalid soil ID

```bash
curl -i http://localhost:8090/api/v1/soil/invalid-id
```

Verified that an invalid UUID returns:

```text
400 Bad Request
```

### Test 7: Non-existent soil record

```bash
curl -i http://localhost:8090/api/v1/soil/00000000-0000-0000-0000-000000000000
```

Verified that a missing record returns:

```text
404 Not Found
```

---

# 9. Farms Module API Tests

The Farms module was tested using the development farm.

### Development farm

| Field | Value |
|---|---|
| Farm name | Demo Farm |
| County | Kisumu |
| Sub-county | Kisumu East |
| Area | 2.5 hectares |
| Latitude | -0.0917 |
| Longitude | 34.7680 |

### Test 1: Retrieve farm

```bash
curl -i http://localhost:8090/api/v1/farms/192f23b1-aedb-4e78-8af7-81c735fb695c
```

Verified that the existing farm returned:

```text
200 OK
```

### Test 2: Invalid farm ID

```bash
curl -i http://localhost:8090/api/v1/farms/invalid-id
```

Verified response:

```text
400 Bad Request
```

### Test 3: Non-existent farm

```bash
curl -i http://localhost:8090/api/v1/farms/00000000-0000-0000-0000-000000000000
```

Verified response:

```text
404 Not Found
```

---

# 10. Crops Module API Tests

The Crops module manages crop records and crop characteristics.

The following development crop records were used:

| Crop | Crop ID |
|---|---|
| Beans | `7c0d0b4b-8764-4010-a714-e10ebdc767d4` |
| Cowpeas | `bd8deeff-4b2b-4015-8e89-2fc2b0e1eeef` |
| Groundnuts | `b15ae5dd-75dc-4225-8f29-fc10fd020222` |
| Maize | `2166cf13-1fdb-489f-9996-86e87c74883e` |
| Millet | `a2c1cc06-4f7f-4f24-a43f-7175f0e43139` |
| Rice | `cb84bbac-c451-4dda-ac5e-022b7a777957` |
| Sorghum | `6703f38c-68c5-4cde-bd74-625088654955` |

### Test 1: Retrieve Beans characteristics

```bash
curl -i http://localhost:8090/api/v1/crops/7c0d0b4b-8764-4010-a714-e10ebdc767d4/characteristics
```

After characteristics were added, the endpoint returned:

```text
200 OK
```

Verified characteristics included:

- Water requirement
- Growth duration
- Nitrogen-fixing status
- Nutrient demand
- Suitable seasons

### Test 2: Retrieve Maize characteristics

```bash
curl -i http://localhost:8090/api/v1/crops/2166cf13-1fdb-489f-9996-86e87c74883e/characteristics
```

Verified response:

```text
200 OK
```

### Test 3: Crop without characteristics

Before adding characteristics for Beans, the endpoint returned:

```text
404 Not Found
```

The response correctly indicated that crop characteristics were not found.

### Test 4: Non-existent crop

```bash
curl -i http://localhost:8090/api/v1/crops/00000000-0000-0000-0000-000000000000/characteristics
```

Verified response:

```text
404 Not Found
```

---

# 11. Earthdata Module

The Earthdata module integrates with the NASA POWER API to retrieve daily climate information.

### Supported climate parameters

- Precipitation
- Average temperature
- Maximum temperature
- Minimum temperature
- Relative humidity

### Endpoint

```text
POST /api/v1/earthdata/climate
```

### Request structure

```json
{
  "farm_id": "192f23b1-aedb-4e78-8af7-81c735fb695c",
  "latitude": -0.0917,
  "longitude": 34.7680,
  "start_date": "2026-01-01",
  "end_date": "2026-01-07"
}
```

### Implementation details

- NASA POWER API client
- Daily climate data retrieval
- Request validation
- HTTP timeout handling
- Climate observation response formatting

**Testing status:** The module was implemented and compiled. A successful live NASA POWER API request has not yet been confirmed in the documented tests.

**Important:** This endpoint retrieves climate observations. It does not currently save them into the PostgreSQL `climate_data` table.

---

# 12. Crop Rotation Module

The Rotation module generates crop-rotation recommendations using a rule-based engine.

### Available endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/api/v1/rotation/` | Generate rotation plans |
| GET | `/api/v1/rotation/:id` | Retrieve rotation plan |
| GET | `/api/v1/rotation/:id/results` | Retrieve rotation results |
| GET | `/api/v1/farms/:id/rotations` | List farm rotation plans |

### Test 1: Generate rotation plans

```bash
curl -i -X POST http://localhost:8090/api/v1/rotation/ \
-H "Content-Type: application/json" \
-d '{
  "farm_id": "192f23b1-aedb-4e78-8af7-81c735fb695c",
  "farmer_priorities": [
    "soil_health",
    "water_efficiency",
    "crop_diversity"
  ],
  "number_of_seasons": 3
}'
```

Verified response:

```text
201 Created
```

The response confirmed:

```json
{
  "message": "Rotation plans generated successfully",
  "success": true
}
```

### Generated rotation strategies

Three strategies were generated:

**1. Water Conscious Rotation**

```text
Cowpeas → Millet → Beans
```

Overall score: 62

**2. Balanced Rotation**

```text
Beans → Cowpeas → Groundnuts
```

Overall score: 60

**3. Soil Restoration Focus**

```text
Beans → Cowpeas → Groundnuts
```

Overall score: 60

Generated plan ID:

```text
**Project:** SoilCycle by Earth2Farm
**Backend language:** Go
**Database:** PostgreSQL 16
**Development API:** `http://localhost:8090/api/v1`9fc7a6c8-3065-49b2-aba5-7aa8cad9db1d
```

### Rotation engine limitations

The current implementation is a development-stage rule-based engine.

- Climate score is currently a placeholder.
- Actual soil measurements are not yet integrated into the rotation scoring.
- Some generated strategies can contain identical crop sequences.
- Recommendations depend on the available crop records and characteristics.
- Agronomic suitability requires local validation.

The scores are provisional decision-support indicators, not verified predictions of yield, profit, or soil improvement.

---

# 13. Development Test Data

The following records were used for development and manual API testing.

### Demo farm

```text
Farm ID:
192f23b1-aedb-4e78-8af7-81c735fb695c

Farm name:
Demo Farm

Location:
Kisumu East, Kisumu County, Kenya
```

### Development soil record

```text
Soil record ID:
0ea0aacf-0655-43ba-9214-aa61e7470b09

Farm ID:
192f23b1-aedb-4e78-8af7-81c735fb695c

pH:
6.8

Organic carbon:
1.35

Nitrogen:
0.20

Phosphorus:
13.0

Potassium:
0.40

Soil texture:
Clay Loam

Data source:
Farmer Entry
```

These are development test values, not verified laboratory measurements.

### Crop characteristics

Seven crop-characteristics records were inserted directly into the development database.

They include:

- Beans
- Cowpeas
- Groundnuts
- Maize
- Millet
- Rice
- Sorghum

These records should be maintained in a database seed script or migration so they can be recreated in another environment.

---

# 14. Testing Summary

| Module | Test | Result |
|---|---|---|
| Docker | Start PostgreSQL container | Passed |
| PostgreSQL | Connect to database | Passed |
| PostgreSQL | Verify database tables | Passed |
| Backend | Go formatting | Passed |
| Backend | `go vet ./...` | Passed |
| Backend | `go build ./...` | Passed |
| Backend | `go test ./...` command | Passed; no unit test files |
| Soil | Create soil record | Passed |
| Soil | List soil records | Passed |
| Soil | Retrieve soil record | Passed |
| Soil | Retrieve soil by farm | Passed |
| Soil | Update soil record | Passed |
| Soil | Invalid UUID handling | Passed |
| Soil | Missing record handling | Passed |
| Farms | Retrieve existing farm | Passed |
| Farms | Invalid farm UUID | Passed |
| Farms | Missing farm | Passed |
| Crops | Retrieve Beans characteristics | Passed |
| Crops | Retrieve Maize characteristics | Passed |
| Crops | Missing characteristics handling | Passed |
| Crops | Missing crop handling | Passed |
| Rotation | Generate rotation plans | Passed |
| Earthdata | Compile and build | Passed |
| Earthdata | Live API request | Pending |
| Rotation | Retrieve generated plan | Pending |
| Rotation | Retrieve rotation results | Pending |
| Rotation | List farm rotation plans | Pending |

---

# 15. Final Development Commands

Before committing backend changes:

```bash
gofmt -w .
go vet ./...
go build ./...
go test ./...
git diff --check
```

Review changes:

```bash
git status
git diff --stat
```

Stage backend files:

```bash
git add README.md main.go go.mod go.sum internal/modules/crops internal/modules/farms internal/modules/earthdata internal/modules/rotation
```

Review staged changes:

```bash
git status
git diff --cached --stat
```

Commit:

```bash
git commit -m "feat: add earthdata and crop rotation modules"
```

Push:

```bash
git push origin feat/backend
```

---

## 16. Future Improvements

1. Add automated Go unit and integration tests.
2. Add reproducible crop-characteristics seed data.
3. Persist NASA POWER climate observations in PostgreSQL.
4. Integrate soil and climate information into rotation scoring.
5. Improve differentiation between rotation strategies.
6. Add authentication and authorization.
7. Connect the backend to the frontend application.
8. Validate crop recommendations with local agricultural experts.
9. Add deployment configuration and production environment settings.
---
**Project:** SoilCycle by Earth2Farm

**Backend language:** Go

**Database:** PostgreSQL 16

**Development API:** `http://localhost:8090/api/v1`
