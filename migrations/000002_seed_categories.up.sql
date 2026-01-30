INSERT INTO categories (name) VALUES ('Eldar')
ON CONFLICT (name) DO NOTHING;
