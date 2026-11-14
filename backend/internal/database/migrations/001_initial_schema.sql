-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. USERS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash TEXT,
    role VARCHAR(30) NOT NULL DEFAULT 'farmer',
    preferred_language VARCHAR(20) DEFAULT 'en',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_user_role
        CHECK (role IN ('farmer', 'admin', 'agronomist'))
);

-- 2. FARMS
CREATE TABLE farms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    county VARCHAR(100) NOT NULL DEFAULT 'Kisumu',
    sub_county VARCHAR(100),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    area_hectares NUMERIC(10,2),
    water_availability VARCHAR(30),
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_latitude
        CHECK (latitude BETWEEN -90 AND 90),

    CONSTRAINT valid_longitude
        CHECK (longitude BETWEEN -180 AND 180),

    CONSTRAINT valid_farm_area
        CHECK (area_hectares IS NULL OR area_hectares > 0)
);

-- 3. CROPS
CREATE TABLE crops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    scientific_name VARCHAR(150),
    crop_family VARCHAR(100),
    category VARCHAR(50),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. CROP CHARACTERISTICS
CREATE TABLE crop_characteristics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_id UUID NOT NULL UNIQUE REFERENCES crops(id) ON DELETE CASCADE,
    water_requirement VARCHAR(30),
    growth_duration_days INTEGER,
    nitrogen_fixing BOOLEAN DEFAULT FALSE,
    nutrient_demand VARCHAR(30),
    suitable_seasons TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_growth_duration
        CHECK (
            growth_duration_days IS NULL
            OR growth_duration_days > 0
        )
);

-- 5. FARM CROPPING HISTORY
CREATE TABLE farm_crops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    crop_id UUID NOT NULL REFERENCES crops(id),
    season VARCHAR(100),
    planting_date DATE,
    harvest_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_crop_dates
        CHECK (
            harvest_date IS NULL
            OR planting_date IS NULL
            OR harvest_date >= planting_date
        )
);

-- 6. SOIL DATA
CREATE TABLE soil_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    ph NUMERIC(4,2),
    organic_carbon NUMERIC(10,3),
    nitrogen NUMERIC(10,3),
    phosphorus NUMERIC(10,3),
    potassium NUMERIC(10,3),
    soil_texture VARCHAR(100),
    data_source VARCHAR(50) NOT NULL,
    measurement_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_soil_ph
        CHECK (ph IS NULL OR ph BETWEEN 0 AND 14)
);

-- 7. CLIMATE DATA
CREATE TABLE climate_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    source VARCHAR(50) NOT NULL,
    observation_date DATE NOT NULL,
    rainfall_mm NUMERIC(10,2),
    temperature_c NUMERIC(6,2),
    humidity_percent NUMERIC(5,2),
    soil_moisture NUMERIC(10,4),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_rainfall
        CHECK (rainfall_mm IS NULL OR rainfall_mm >= 0),

    CONSTRAINT valid_humidity
        CHECK (
            humidity_percent IS NULL
            OR humidity_percent BETWEEN 0 AND 100
        )
);

-- 8. ROTATION PLANS
CREATE TABLE rotation_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    farmer_priorities JSONB NOT NULL DEFAULT '[]'::jsonb,
    input_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    number_of_seasons INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_number_of_seasons
        CHECK (number_of_seasons BETWEEN 1 AND 10)
);

-- 9. ROTATION RESULTS
CREATE TABLE rotation_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES rotation_plans(id) ON DELETE CASCADE,
    plan_name VARCHAR(150) NOT NULL,
    crop_sequence JSONB NOT NULL,
    scores JSONB NOT NULL DEFAULT '{}'::jsonb,
    explanation JSONB NOT NULL DEFAULT '[]'::jsonb,
    overall_score NUMERIC(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_overall_score
        CHECK (
            overall_score IS NULL
            OR overall_score BETWEEN 0 AND 100
        )
);

-- INDEXES

CREATE INDEX idx_farms_user_id
ON farms(user_id);

CREATE INDEX idx_farm_crops_farm_id
ON farm_crops(farm_id);

CREATE INDEX idx_farm_crops_crop_id
ON farm_crops(crop_id);

CREATE INDEX idx_soil_data_farm_id
ON soil_data(farm_id);

CREATE INDEX idx_climate_data_farm_date
ON climate_data(farm_id, observation_date);

CREATE INDEX idx_rotation_plans_farm_id
ON rotation_plans(farm_id);

CREATE INDEX idx_rotation_results_plan_id
ON rotation_results(plan_id);
