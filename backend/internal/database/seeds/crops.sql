INSERT INTO crops
(name, scientific_name, crop_family, category, description)
VALUES
(
    'Maize',
    'Zea mays',
    'Poaceae',
    'cereal',
    'Important staple cereal crop.'
),
(
    'Beans',
    'Phaseolus vulgaris',
    'Fabaceae',
    'legume',
    'Common food legume with nitrogen-fixing potential.'
),
(
    'Sorghum',
    'Sorghum bicolor',
    'Poaceae',
    'cereal',
    'Cereal crop with relatively good drought tolerance.'
),
(
    'Cowpeas',
    'Vigna unguiculata',
    'Fabaceae',
    'legume',
    'Food legume with nitrogen-fixing potential.'
),
(
    'Groundnuts',
    'Arachis hypogaea',
    'Fabaceae',
    'legume',
    'Oilseed legume crop.'
),
(
    'Millet',
    'Pennisetum glaucum',
    'Poaceae',
    'cereal',
    'Cereal crop adapted to relatively dry conditions.'
)
ON CONFLICT (name) DO NOTHING;
