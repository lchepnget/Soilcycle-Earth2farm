BEGIN;

INSERT INTO crops (
    name,
    scientific_name,
    crop_family,
    category,
    description
)
VALUES
    (
        'Maize',
        'Zea mays',
        'Poaceae',
        'cereal',
        'Cereal grain used for food and feed; maturity and water needs vary by variety and location.'
    ),
    (
        'Beans',
        'Phaseolus vulgaris',
        'Fabaceae',
        'legume',
        'Common bean food legume; nodulation and nitrogen fixation depend on soil conditions and compatible rhizobia.'
    ),
    (
        'Sorghum',
        'Sorghum bicolor',
        'Poaceae',
        'cereal',
        'Cereal grain used for food and feed; adaptation and maturity vary by cultivar and site.'
    ),
    (
        'Cowpeas',
        'Vigna unguiculata',
        'Fabaceae',
        'legume',
        'Food legume used for grain and leaves; growth habit and maturity vary by cultivar and site.'
    ),
    (
        'Groundnuts',
        'Arachis hypogaea',
        'Fabaceae',
        'legume',
        'Legume grown for edible seeds and oil; crop performance depends on variety and local conditions.'
    ),
    (
        'Millet',
        'Pennisetum glaucum',
        'Poaceae',
        'cereal',
        'Pearl millet cereal grown for grain; maturity varies by cultivar and environment.'
    ),
    (
        'Rice',
        'Oryza sativa',
        'Poaceae',
        'cereal',
        'Rice grain crop; cultivation system, maturity and water supply vary by location.'
    ),
    (
        'Green grams',
        'Vigna radiata',
        'Fabaceae',
        'legume',
        'Mung bean food legume; maturity and water needs vary by cultivar and site.'
    ),
    (
        'Pigeon peas',
        'Cajanus cajan',
        'Fabaceae',
        'legume',
        'Food legume with short- and long-duration types; maturity varies substantially by cultivar.'
    ),
    (
        'Soybeans',
        'Glycine max',
        'Fabaceae',
        'legume',
        'Oil and protein legume; maturity and nodulation depend on cultivar and local conditions.'
    ),
    (
        'Cassava',
        'Manihot esculenta',
        'Euphorbiaceae',
        'root_tuber',
        'Root crop harvested at variable ages; crop duration depends strongly on variety and intended use.'
    ),
    (
        'Sweet potatoes',
        'Ipomoea batatas',
        'Convolvulaceae',
        'root_tuber',
        'Root crop propagated vegetatively; harvest timing varies by variety and market.'
    ),
    (
        'Tomatoes',
        'Solanum lycopersicum',
        'Solanaceae',
        'vegetable',
        'Fruit vegetable; maturity and water needs depend on cultivar and production system.'
    ),
    (
        'Kale',
        'Brassica oleracea var. acephala',
        'Brassicaceae',
        'vegetable',
        'Leafy vegetable often harvested over multiple pickings; first harvest and crop duration vary.'
    ),
    (
        'Cabbage',
        'Brassica oleracea var. capitata',
        'Brassicaceae',
        'vegetable',
        'Head-forming vegetable; crop cycle varies by cultivar, climate and management.'
    )
ON CONFLICT (name) DO NOTHING;

-- Durations are indicative planning values, not universal maturity dates.
-- Suitable seasons are left unset because planting calendars vary by county and production system.
INSERT INTO crop_characteristics (
    crop_id,
    water_requirement,
    growth_duration_days,
    nitrogen_fixing,
    nutrient_demand,
    suitable_seasons
)
SELECT
    c.id,
    catalogue.water_requirement,
    catalogue.growth_duration_days,
    catalogue.nitrogen_fixing,
    catalogue.nutrient_demand,
    NULL::TEXT[]
FROM (
    VALUES
        ('Maize', 'medium', 120, FALSE, 'high'),
        ('Beans', 'medium', 90, TRUE, 'medium'),
        ('Sorghum', 'low', 110, FALSE, 'medium'),
        ('Cowpeas', 'low', 75, TRUE, 'low'),
        ('Groundnuts', 'medium', 120, TRUE, 'medium'),
        ('Millet', 'low', 90, FALSE, 'low'),
        ('Rice', 'high', 120, FALSE, 'high'),
        ('Green grams', 'low', 75, TRUE, 'low'),
        ('Pigeon peas', 'low', 180, TRUE, 'medium'),
        ('Soybeans', 'medium', 120, TRUE, 'medium'),
        ('Cassava', 'low', 270, FALSE, 'low'),
        ('Sweet potatoes', 'medium', 120, FALSE, 'medium'),
        ('Tomatoes', 'high', 120, FALSE, 'high'),
        ('Kale', 'medium', 90, FALSE, 'medium'),
        ('Cabbage', 'medium', 100, FALSE, 'high')
) AS catalogue(
    name,
    water_requirement,
    growth_duration_days,
    nitrogen_fixing,
    nutrient_demand
)
JOIN crops c ON c.name = catalogue.name
ON CONFLICT (crop_id) DO NOTHING;

COMMIT;
