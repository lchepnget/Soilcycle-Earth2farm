# SoilCycle by Earth2Farm
## Project Overview

### 1. Introduction

SoilCycle by Earth2Farm is a digital agricultural decision-support platform designed to help farmers make informed decisions about crop rotation, soil health, and climate-resilient farming.

The platform brings together NASA Earth observation data, soil information, crop characteristics, and farmer priorities to provide practical insights into different crop rotation options.

The project aims to make agricultural and environmental information easier to understand and apply, particularly for smallholder farmers.

### 2. Problem Statement

Smallholder farmers face several challenges that affect agricultural productivity and long-term soil health.

These include:

- Declining soil fertility due to continuous cultivation and unsuitable cropping practices.
- Limited access to reliable soil information and agricultural advisory services.
- Changing rainfall patterns and other climate-related challenges.
- Limited knowledge of suitable crop rotation practices.
- Difficulty accessing and interpreting environmental and agricultural data.
- Limited access to decision-support tools that consider local farming conditions and farmer priorities.

Many farmers make crop selection and rotation decisions without sufficient information about their soil conditions, climate patterns, and the long-term effects of their choices.

### 3. Proposed Solution

SoilCycle by Earth2Farm will provide a web-based platform that combines Earth observation information, available soil data, crop characteristics, and farmer-provided information.

Farmers will be able to provide information such as:

- Farm location.
- Previous crops grown.
- Available soil test results.
- Water availability.
- Farming priorities, such as food production, income generation, and soil improvement.

The platform will use this information alongside relevant environmental datasets to compare possible crop rotation scenarios.

It will explain the potential advantages, disadvantages, and limitations of each option to support informed decision-making.

### 4. Project Objectives

The main objective is to develop an accessible digital platform that supports sustainable crop rotation and improved soil management.

Specific objectives are to:

1. Integrate relevant NASA Earth observation datasets into an agricultural decision-support system.
2. Incorporate available soil information and crop characteristics.
3. Help farmers explore and compare crop rotation options.
4. Promote practices that support long-term soil health.
5. Support climate-informed agricultural planning.
6. Make environmental and agricultural information understandable to farmers and agricultural stakeholders.

### 5. Target Users

The platform will serve the following groups:

**Primary users**

- Smallholder farmers seeking practical guidance on crop rotation and soil management.

**Secondary users**

- Agricultural extension officers supporting farmers.
- Farmer groups and cooperatives.
- Agricultural organizations and development programmes.

The platform will be designed with farmers' needs and accessibility in mind.

### 6. Core Features

The initial version will focus on the following features:

1. **Farm Profile:** Capture basic farm information, location, and farming conditions.
2. **Crop History:** Allow users to record crops previously grown on a farm.
3. **Soil Information:** Incorporate available soil test results and relevant soil datasets.
4. **Earth Observation Integration:** Retrieve relevant environmental information from NASA-supported data sources.
5. **Crop Rotation Recommendations:** Generate possible crop rotation options based on available information.
6. **Scenario Comparison:** Allow users to compare different rotation options.
7. **Decision Explanations:** Explain why a rotation option may be suitable, including its benefits, trade-offs, and limitations.

### 7. Earth Observation and Data Sources

The project intends to explore the following data sources:

- **NASA POWER:** Weather and solar-related environmental information.
- **NASA GPM:** Satellite-based precipitation estimates.
- **NASA SMAP:** Soil moisture observations.
- **NASA MODIS:** Vegetation and land-surface observations.
- **SoilGrids (ISRIC):** Complementary soil property information.
- **Local Soil Test Results:** Field-level soil information where available.

These datasets will be used according to their spatial resolution, temporal coverage, and documented limitations.

Satellite observations will not be presented as exact measurements of an individual farm's soil conditions. Where field-level information is unavailable, the platform will clearly communicate uncertainty.

### 8. Project Scope

The initial development will focus on building a functional minimum viable product (MVP).

The MVP will aim to:

- Capture basic farm information.
- Record previous crops.
- Integrate selected environmental datasets.
- Accept available soil information.
- Compare selected crop rotation scenarios.
- Present understandable recommendations and explanations.

The initial pilot may focus on Kisumu County, Kenya, while the system will be designed with the potential to support farmers in other parts of Kenya and, eventually, Africa.

### 9. Expected Impact

SoilCycle by Earth2Farm aims to:

- Promote better crop rotation practices.
- Support improved soil management.
- Increase awareness of climate-related agricultural conditions.
- Help farmers make more informed cropping decisions.
- Support extension officers and agricultural organizations with accessible information.
- Contribute to more sustainable and resilient farming systems.

The project will evaluate its impact through practical testing, user feedback, and the usefulness of its recommendations.

### 10. Technology Direction

The proposed technology stack includes:

- **Backend:** Go.
- **Frontend:** HTML, CSS, and JavaScript, with the option to adopt a frontend framework if agreed by the team.
- **Database:** PostgreSQL.
- **Geospatial capabilities:** PostGIS, if required.
- **External data:** NASA and complementary soil-data services.
- **Development and deployment:** Git, GitHub, Docker, and automated testing tools.

The final architecture and technology decisions will be confirmed by the development team.

### 11. Project Vision

Our vision is to make Earth observation and agricultural data useful to farmers by turning complex environmental information into understandable, practical farming decisions.

SoilCycle by Earth2Farm seeks to contribute to a future where farmers can make better-informed decisions, maintain healthier soils, and build more resilient agricultural systems.

**Tagline:** Smarter rotations. Healthier soils. Resilient farming.