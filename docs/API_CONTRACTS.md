# API Contracts
## SoilCycle by Earth2Farm

### 1. Purpose

This document describes the proposed API contracts for SoilCycle by Earth2Farm.

The API will allow the frontend application to communicate with the Go backend, submit farm information, retrieve environmental data, and request crop rotation recommendations.

These contracts are initial proposals and may be updated following team discussions and implementation decisions.

### 2. API Design

- API style: REST.
- Data format: JSON.
- Backend language: Go.
- Base URL (development): `http://localhost:8080/api/v1`
- Communication protocol: HTTP/HTTPS.

### 3. Proposed Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/health` | Check backend availability |
| POST | `/farms` | Register a farm |
| GET | `/farms/{id}` | Retrieve farm information |
| PUT | `/farms/{id}` | Update farm information |
| POST | `/farms/{id}/crops` | Record crop history |
| GET | `/farms/{id}/crops` | Retrieve crop history |
| POST | `/farms/{id}/soil-data` | Submit available soil information |
| GET | `/farms/{id}/environment` | Retrieve environmental information |
| POST | `/farms/{id}/recommendations` | Generate crop rotation recommendations |
| POST | `/farms/{id}/scenarios/compare` | Compare crop rotation scenarios |

### 4. Request and Response Examples

#### 4.1 Health Check

**Request**

`GET /api/v1/health`

**Example response**

```json
{
  "status": "ok",
  "service": "SoilCycle by Earth2Farm"
}
```

#### 4.2 Register a Farm

**Request**

`POST /api/v1/farms`

**Example request body**

```json
{
  "farm_name": "Demo Farm",
  "location": {
    "latitude": -0.1022,
    "longitude": 34.7617
  },
  "area_hectares": 2.5,
  "water_availability": "rain-fed"
}
```

**Example response**

```json
{
  "message": "Farm registered successfully",
  "farm_id": "farm-001"
}
```

#### 4.3 Record Crop History

**Request**

`POST /api/v1/farms/{id}/crops`

**Example request body**

```json
{
  "crop_name": "Maize",
  "season": "2026-long-rains",
  "year": 2026
}
```

**Example response**

```json
{
  "message": "Crop history recorded successfully"
}
```

#### 4.4 Submit Soil Information

**Request**

`POST /api/v1/farms/{id}/soil-data`

**Example request body**

```json
{
  "ph": 6.2,
  "organic_carbon": 1.8,
  "nitrogen": 0.15,
  "source": "laboratory-soil-test"
}
```

**Example response**

```json
{
  "message": "Soil information received successfully"
}
```

The example measurements are illustrative, not actual farm results or recommended soil values. Actual fields and units will be agreed during database and backend design.

#### 4.5 Retrieve Environmental Information

**Request**

`GET /api/v1/farms/{id}/environment`

**Example response**

```json
{
  "farm_id": "farm-001",
  "data_sources": [
    "NASA POWER",
    "NASA GPM",
    "NASA SMAP",
    "NASA MODIS"
  ],
  "message": "Environmental information retrieved successfully"
}
```

Actual response fields will depend on the datasets and variables integrated.

#### 4.6 Generate Crop Rotation Recommendations

**Request**

`POST /api/v1/farms/{id}/recommendations`

**Example request body**

```json
{
  "priority": "soil_health",
  "planning_seasons": 3
}
```

**Example response**

```json
{
  "farm_id": "farm-001",
  "recommendations": [
    {
      "rotation": [
        "Maize",
        "Beans",
        "Groundnuts"
      ],
      "reason": "This option includes legumes that may contribute to soil fertility when appropriately managed."
    }
  ]
}
```

Recommendations will depend on crop suitability, available environmental information, previous crops, and farmer priorities.

#### 4.7 Compare Rotation Scenarios

**Request**

`POST /api/v1/farms/{id}/scenarios/compare`

**Example request body**

```json
{
  "scenarios": [
    {
      "name": "Scenario A",
      "crops": ["Maize", "Beans", "Groundnuts"]
    },
    {
      "name": "Scenario B",
      "crops": ["Sorghum", "Cowpeas", "Maize"]
    }
  ]
}
```

**Example response**

```json
{
  "comparison": [
    {
      "scenario": "Scenario A",
      "assessment": "Includes legume crops that may support soil fertility."
    },
    {
      "scenario": "Scenario B",
      "assessment": "Provides an alternative rotation for evaluation under local conditions."
    }
  ]
}
```

The comparison engine will use documented agricultural rules and available data. It will not present unvalidated estimates as scientifically proven results.

### 5. Standard Error Responses

The API should use consistent HTTP status codes.

| Status code | Meaning |
|---|---|
| 200 | Request successful |
| 201 | Resource created |
| 400 | Invalid request |
| 404 | Resource not found |
| 422 | Request understood but cannot be processed |
| 500 | Internal server error |
| 503 | External data service unavailable |

**Example error response**

```json
{
  "error": "Invalid request",
  "message": "Farm location is required."
}
```

### 6. Validation and Data Quality

The backend should validate:

- Required farm information.
- Valid geographic coordinates.
- Positive farm area.
- Valid crop and season information.
- Soil measurements and their units.
- Supported environmental data parameters.

Missing or unavailable environmental information should be communicated clearly rather than replaced with invented values.

### 7. API Security

The implementation should consider:

- Input validation.
- Secure handling of user information.
- Authentication and authorization where required.
- HTTPS in production.
- Protection against unauthorized access.
- Safe handling of external API credentials.

### 8. Development Notes

1. These endpoints are proposed and require team review.
2. Final request and response fields will be aligned with the database schema.
3. NASA data integration will be handled by the backend.
4. The recommendation engine will initially use transparent, documented rules.
5. API versioning will use `/api/v1`.
6. The team should update this document whenever API contracts change.

### 9. Future Improvements

Potential future API capabilities include:

- User authentication and account management.
- Farmer group management.
- Extension officer access.
- Historical environmental data retrieval.
- Recommendation history.
- Exporting reports.

**Document status:** Initial proposal — pending team review.