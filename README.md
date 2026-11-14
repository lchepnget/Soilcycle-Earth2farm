# SoilCycle by Earth2Farm

### NASA-Powered Climate-Smart Crop Rotation and Soil Health Decision-Support Platform

**Building climate-resilient agriculture through Earth observation, soil intelligence, and smarter farming decisions.**

---

## 1. Project Overview

SoilCycle by Earth2Farm is a technology-driven agricultural decision-support platform designed to help smallholder farmers make informed crop rotation and soil management decisions using NASA Earth observations, local soil information, crop characteristics, and farmer priorities.

The project aims to address challenges associated with changing rainfall patterns, rising temperatures, water availability, declining soil health, and unsustainable cropping practices.

Starting with smallholder farmers in Kisumu County, Kenya, the platform is designed with the potential to expand across Kenya and other African agricultural regions.

## 2. Problem Statement

Smallholder farmers face increasing agricultural uncertainty due to climate variability, changing rainfall patterns, water shortages, and declining soil fertility.

Many farmers lack accessible tools that combine environmental observations, soil information, crop characteristics, and their own farming history to guide crop rotation decisions.

As a result, farmers may continue using cropping practices that do not adequately address soil health, water management, and climate resilience.

There is a need for an accessible, data-driven decision-support solution that helps farmers explore suitable crop rotation strategies and understand their potential long-term effects on agricultural sustainability.

## 3. Proposed Solution

SoilCycle by Earth2Farm will provide a web-based platform that integrates NASA Earth observations with local soil information, crop characteristics, and farmer inputs.

The platform will analyse relevant environmental and agricultural information to generate alternative crop rotation strategies.

Farmers will be able to compare possible rotations, understand their potential soil-health benefits, assess water-related considerations, and explore strategies aligned with their farming priorities.

Rather than providing generic crop recommendations, SoilCycle will focus on explainable, evidence-informed agricultural decision support.

## 4. Project Objectives

- Promote sustainable crop rotation practices.
- Support soil health management and conservation.
- Help farmers understand changing environmental conditions.
- Integrate NASA Earth observations into practical agricultural decision-making.
- Provide accessible agricultural insights for smallholder farmers.
- Support climate-resilient and sustainable farming systems.
- Develop a scalable solution that can serve agricultural communities across Africa.

## 5. Core Features

### 5.1 Farmer and Farm Profile
- Farm location selection.
- Previous crop history.
- Available soil information.
- Water availability.
- Farmer priorities and agricultural goals.

### 5.2 Earth Observation Intelligence
- Integration of NASA environmental datasets.
- Rainfall and temperature analysis.
- Regional soil moisture indicators.
- Vegetation and land-surface condition monitoring.
- Environmental trend visualisation.

### 5.3 Soil Health Insights
- Interpretation of available soil information.
- Soil-health management guidance.
- Crop nutrient-demand considerations.
- Soil conservation recommendations.

### 5.4 Crop Rotation Simulator
- Alternative crop rotation generation.
- Crop sequence comparisons.
- Crop characteristics analysis.
- Potential soil-health benefits and limitations.
- Water and climate suitability considerations.

### 5.5 Farmer Decision Dashboard
- Visual comparison of rotation strategies.
- Environmental information dashboards.
- Explainable suitability assessments.
- Farmer-oriented recommendations.
- Seasonal planning support.

## 6. NASA Earth Observation Data Sources

The platform will explore the following NASA data sources:

| Dataset | Intended application |
|---|---|
| NASA POWER | Historical and near-real-time meteorological and solar-related data |
| NASA GPM | Satellite-based precipitation observations |
| NASA SMAP | Regional soil moisture information |
| NASA MODIS | Vegetation and land-surface observations |

Additional complementary data sources may include:

- SoilGrids for soil property information.
- Local soil-testing results.
- Relevant Kenyan agricultural and climate datasets.
- Documented crop characteristics and agronomic knowledge.

**Data integrity:** Satellite observations will be interpreted according to their spatial and temporal limitations. Regional estimates will not be presented as direct measurements of individual farms. Historical observations will not be presented as guaranteed future forecasts.

## 7. Technology Stack

| Component | Technology |
|---|---|
| Backend | Go (Golang) |
| Frontend | HTML, CSS, JavaScript / suitable frontend framework |
| Database | PostgreSQL |
| Database management | SQL migrations |
| Geospatial capabilities | PostGIS (where required) |
| External data integration | NASA APIs and relevant public datasets |
| Containerization | Docker |
| Version control | Git and GitHub |
| Continuous integration | GitHub Actions |

## 8. Repository Structure

```text
soilcycle-earth2farm/
│
├── backend/
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── tests/
│   ├── docs/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── frontend/
│   ├── public/
│   ├── src/
│   ├── Dockerfile
│   ├── package.json
│   └── README.md
│
├── database/
│   ├── migrations/
│   │   ├── 001_create_users.sql
│   │   ├── 002_create_farms.sql
│   │   ├── 003_create_crops.sql
│   │   ├── 004_create_soil_data.sql
│   │   └── 005_create_rotations.sql
│   │
│   └── seeds/
│       └── crops.sql
│
├── docs/
│   ├── PROJECT_OVERVIEW.md
│   ├── ARCHITECTURE.md
│   ├── API_CONTRACTS.md
│   ├── DATABASE_SCHEMA.md
│   ├── NASA_DATA_SOURCES.md
│   ├── TESTING_PLAN.md
│   └── CONTRIBUTING.md
│
├── scripts/
│   ├── setup.sh
│   └── dev.sh
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── .gitignore
├── docker-compose.yml
└── README.md
```

## 9. Target Users

- Smallholder farmers.
- Agricultural extension officers.
- Farmer groups and cooperatives.
- Agricultural researchers.
- Community-based agricultural organisations.
- Agricultural development practitioners.

## 10. Expected Impact

SoilCycle by Earth2Farm aims to:

- Improve access to agricultural and environmental information.
- Encourage informed crop rotation decisions.
- Promote sustainable soil management.
- Support climate-resilient agricultural practices.
- Strengthen farmers' ability to understand environmental risks.
- Contribute to sustainable food production and agricultural livelihoods.

The long-term vision is to make accessible Earth observation-powered agricultural decision support available to farming communities across Africa.

## 11. Team and Collaboration

We are a five-member, technology-driven team focused on building clean, impactful, and structurally sound digital solutions.

Our collective interests include software engineering, full-stack web development, backend systems, data integration, and transforming complex datasets into accessible applications.

**Project Management:** Project coordination, planning, documentation, and team collaboration.

Team member profiles and GitHub links will be added as the team finalises its responsibilities.

## 12. Development Status

**Current stage:** Project initiation and repository setup.

Development will proceed through:
1. Project planning and requirements definition.
2. System architecture and database design.
3. NASA data exploration and integration.
4. Backend development.
5. Frontend development.
6. Crop rotation decision engine.
7. Integration and testing.
8. Deployment and demonstration.

## 13. Contribution Guidelines

We welcome organised collaboration from all project team members.

- Work on assigned tasks.
- Create feature branches for individual contributions.
- Use meaningful commit messages.
- Keep documentation updated.
- Test changes before submitting them.
- Submit pull requests for review.
- Avoid pushing unreviewed changes directly to the main branch.

See `docs/CONTRIBUTING.md` for detailed contribution guidelines.

## 14. Project Vision

**From Earth observations to better farming decisions.**

SoilCycle by Earth2Farm envisions a future where farmers can use environmental intelligence and agricultural knowledge to make informed decisions, protect their soil, adapt to changing climate conditions, and build more resilient farming systems.

---

**SoilCycle by Earth2Farm**

*Smarter rotations. Healthier soils. Resilient farming.* 🌍🌱