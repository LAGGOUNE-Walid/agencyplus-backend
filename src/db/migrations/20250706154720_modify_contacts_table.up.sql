ALTER TABLE contacts ADD COLUMN preferred_building_types TEXT; -- comma-separated or JSON
ALTER TABLE contacts ADD COLUMN preferred_features TEXT; -- JSON
ALTER TABLE contacts ADD COLUMN min_rooms INT;
ALTER TABLE contacts ADD COLUMN max_rooms INT;
ALTER TABLE contacts ADD COLUMN min_surface DECIMAL(10, 2);
ALTER TABLE contacts ADD COLUMN max_surface DECIMAL(10, 2);
ALTER TABLE contacts ADD COLUMN furnished BOOLEAN;
ALTER TABLE contacts ADD COLUMN acceptable_payment_type VARCHAR(100);
ALTER TABLE contacts ADD COLUMN max_year_built INT;
