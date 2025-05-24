CREATE TABLE IF NOT EXISTS buildings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id BIGINT NOT NULL,
    location TEXT,
    title VARCHAR(255),
    wilaya VARCHAR(100),
    daira VARCHAR(100),
    building_type VARCHAR(100),
    is_promotion_building BOOLEAN DEFAULT FALSE,
    is_residency BOOLEAN DEFAULT FALSE,
    status VARCHAR(100),
    price INT,
    surface_total DECIMAL(10, 2),
    surface_built DECIMAL(10, 2),
    rooms INT,
    bathrooms INT,
    floors_total INT,
    parking_spaces INT,
    is_by_the_sea BOOLEAN DEFAULT FALSE,
    has_water BOOLEAN DEFAULT FALSE,
    has_electricity BOOLEAN DEFAULT FALSE,
    has_gas BOOLEAN DEFAULT FALSE,
    has_internet BOOLEAN DEFAULT FALSE,
    has_garden BOOLEAN DEFAULT FALSE,
    has_pool BOOLEAN DEFAULT FALSE,
    has_elevator BOOLEAN DEFAULT FALSE,
    has_central_heating BOOLEAN DEFAULT FALSE,
    has_water_tank BOOLEAN DEFAULT FALSE,
    has_air_conditioner BOOLEAN DEFAULT FALSE,
    has_equipped_kitchen BOOLEAN DEFAULT FALSE,
    has_terrace BOOLEAN DEFAULT FALSE,
    has_notarial_deed BOOLEAN DEFAULT FALSE, -- Acte notarié
    has_land_booklet BOOLEAN DEFAULT FALSE, -- Livret foncier
    has_act_in_joint_ownership BOOLEAN DEFAULT FALSE, -- Acte dans l'indivision
    has_certificate_of_conformity BOOLEAN DEFAULT FALSE, -- Certificat de conformité
    has_decision BOOLEAN DEFAULT FALSE,
    has_concession BOOLEAN DEFAULT FALSE,
    has_stamped_paper BOOLEAN DEFAULT FALSE, -- Papier timbré
    has_building_permit BOOLEAN DEFAULT FALSE, -- Permis de construire
    has_off_plan_sales_contract BOOLEAN DEFAULT FALSE, -- Contrat vente sur plan
    building_finished_type VARCHAR(50), -- fini, semi fini, carcasse
    acceptable_payment_type VARCHAR(100), -- tranché, credit bank
    furnished BOOLEAN DEFAULT FALSE,
    year_built YEAR,
    description TEXT,
    shareable_link TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
