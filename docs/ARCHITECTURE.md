# SoilCycle by Earth2Farm
# System Architecture

## 1. Introduction

This document describes the proposed system architecture for SoilCycle by Earth2Farm.

The platform will integrate agricultural information, soil data, crop characteristics, and NASA Earth observation datasets to support crop rotation planning and soil-health decision-making.

The architecture is designed to be modular, scalable, and maintainable, allowing the development team to build and improve individual components.

This document describes the initial architecture proposal. Specific implementation decisions may be updated as the team progresses.

## 2. Architecture Overview

The system will consist of five main components:

1. Frontend application.
2. Backend API.
3. Database.
4. External Earth observation and soil data services.
5. Crop rotation decision engine.

These components will communicate to collect farm information, retrieve relevant environmental data, process available information, and present understandable recommendations to users.

## 3. System Architecture Diagram

```text
                 USERS
       Farmers | Extension Officers
       Farmer Groups | Organizations
                    |
                    v
             FRONTEND
       Web Application Interface
                    |
                    v
              GO BACKEND
                REST API
                    |
          +---------+----------+
          |         |          |
          v         v          v
       DATABASE   DATA       ROTATION
       PostgreSQL SERVICES   ENGINE
          |         |          |
          |         v          |
          |     NASA APIs      |
          |     POWER          |
          |     GPM            |
          |     SMAP           |
          |     MODIS          |
          |                    |
          |     SoilGrids      |
          |                    |
          +---------+----------+
                    |
                    v
             RECOMMENDATIONS
                    |
                    v
                FRONTEND
```

## 4. Frontend Layer

The frontend will provide the interface through which users interact with the platform.

### Responsibilities

- User-friendly farm information forms.
- Farm location input.
- Crop history recording.
- Soil information submission.
- Displaying environmental information.
- Presenting crop rotation scenarios.
- Comparing recommendations.
- Displaying explanations and limitations.

### Proposed Technologies

- HTML
- CSS
- JavaScript

A frontend framework may be considered if agreed upon by the team.

## 5. Backend Layer

The backend will be developed using Go.

It will serve as the central processing layer connecting the frontend, database, external data services, and crop rotation engine.

### Responsibilities

- Receive frontend requests.
- Validate submitted farm information.
- Retrieve stored farm and crop information.
- Communicate with external environmental data services.
- Process available information.
- Run crop rotation evaluation logic.
- Return structured responses to the frontend.

### Proposed Structure

```text
backend/
├── cmd/
├── internal/
├── pkg/
├── tests/
├── docs/
├── Dockerfile
├── go.mod
└── README.md
```

The detailed responsibilities of individual packages will be agreed upon during implementation.

## 6. Database Layer

PostgreSQL is the proposed database.

The database will store information required for farm profiles, crop histories, soil information, and crop rotation scenarios.

### Proposed Data Categories

- User information.
- Farm profiles and locations.
- Crop records.
- Soil information.
- Environmental data references or cached observations.
- Crop rotation scenarios.
- Recommendation records.

PostGIS may be introduced if the project requires advanced geospatial queries.

The final database schema will be documented separately in DATABASE_SCHEMA.md.

## 7. External Data Integration

The platform intends to integrate relevant external environmental datasets.

### NASA POWER

Provides meteorological and solar-related information that can support agricultural climate analysis.

### NASA GPM

Provides satellite-based precipitation estimates that can help assess rainfall patterns.

### NASA SMAP

Provides soil moisture observations useful for understanding broader regional moisture conditions.

### NASA MODIS

Provides vegetation and land-surface observations that can contribute to environmental monitoring.

### SoilGrids

Provides complementary soil property information from ISRIC.

### Local Soil Test Results

Where available, users may provide local soil test information to improve the relevance of farm-level assessments.

External data will be retrieved through suitable services or APIs, subject to availability and access requirements.

The system will communicate the resolution, date, coverage, and limitations of the data used.

## 8. Crop Rotation Decision Engine

The crop rotation engine will be responsible for evaluating possible crop rotation scenarios.

The initial version will use transparent, rule-based logic rather than relying on unvalidated machine-learning predictions.

### Inputs

- Previous crops.
- Available soil information.
- Farm location.
- Environmental conditions.
- Crop characteristics.
- Farmer priorities.

### Processing

The engine will compare available information against documented crop rotation principles and suitability rules.

### Outputs

- Suggested crop rotation options.
- Explanation of why an option may be suitable.
- Potential soil-health considerations.
- Climate-related considerations.
- Risks and trade-offs.
- Information gaps and limitations.

Recommendations will be decision-support suggestions, not guarantees of yield or improved soil fertility.

## 9. System Data Flow

The proposed workflow is:

1. A user opens the web application.
2. The user enters farm details and available agricultural information.
3. The frontend sends the information to the Go backend.
4. The backend validates the submitted information.
5. The backend retrieves relevant stored information.
6. The backend obtains relevant environmental data from external services where available.
7. The crop rotation engine evaluates possible scenarios.
8. The backend returns the results.
9. The frontend displays recommendations, comparisons, explanations, and limitations.

## 10. Security and Reliability Considerations

The system should consider:

- Input validation.
- Secure handling of user information.
- Appropriate authentication and authorization.
- Protection of database credentials.
- Safe handling of external API failures.
- Clear error messages.
- Logging and monitoring.
- Data backup and recovery.

Authentication and authorization requirements will be finalized as the team defines user roles.

## 11. Scalability and Future Improvements

The architecture should allow future improvements such as:

- Additional satellite and agricultural datasets.
- More crop varieties and rotation rules.
- Advanced geospatial analysis.
- Mobile-friendly access.
- Multilingual user interfaces.
- Offline or low-connectivity support.
- More advanced analytical models after appropriate validation.

The platform will initially focus on a practical MVP while allowing future expansion across Kenya and Africa.

## 12. Architecture Decisions Pending Team Agreement

The following decisions require confirmation:

- Final frontend framework.
- Database deployment approach.
- Authentication method.
- Exact NASA data access methods.
- Data caching strategy.
- Geospatial processing requirements.
- Crop rotation evaluation methodology.
- Hosting and deployment environment.

These decisions should be discussed and recorded by the development team before implementation.

## 13. Conclusion

SoilCycle by Earth2Farm will use a modular architecture to connect farm information, environmental observations, soil data, and crop rotation evaluation.

The proposed system aims to transform complex agricultural and environmental information into understandable decision-support insights that help farmers and agricultural stakeholders make better-informed farming decisions.