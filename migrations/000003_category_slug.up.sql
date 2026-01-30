ALTER TABLE categories ADD COLUMN slug VARCHAR(100) UNIQUE;
UPDATE categories SET slug = 'eldar' WHERE name = 'Eldar';
UPDATE categories SET slug = LOWER(REPLACE(TRIM(name), ' ', '-')) WHERE slug IS NULL OR slug = '';
ALTER TABLE categories ALTER COLUMN slug SET NOT NULL;
