CREATE TABLE IF NOT EXISTS contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    fullname TEXT NOT NULL,                  -- Full name
    phone TEXT,                              -- Phone number
    email TEXT,                              -- Email address
    wilaya TEXT,                             -- Wilaya (region) looking for
    daira TEXT,                              -- Daira (sub-region) looking for
    client_type TEXT,                        -- Client type (Acheteur, Locataire, etc.)
    preferred_location_type TEXT,            -- Urban, rural, etc.
    house_finishing TEXT,                    -- Finishing type (e.g. semi-fini, fini)
    renting_floor_looking_for TEXT,          -- Preferred floor if renting
    is_married BOOLEAN DEFAULT 0,            -- Marital status
    preferred_building_types TEXT,           -- JSON or comma-separated (e.g. ['Appartement', 'Villa'])
    preferred_features TEXT,                 -- JSON list of features (e.g. ['Ascenseur', 'Garage'])
    min_rooms INT,                           -- Minimum number of rooms
    max_rooms INT,                           -- Maximum number of rooms
    min_budget INTEGER,
    max_budget INTEGER,
    min_surface DECIMAL(10, 2),              -- Minimum surface area in m²
    max_surface DECIMAL(10, 2),              -- Maximum surface area in m²
    furnished BOOLEAN,                       -- Is the client looking for a furnished property?
    acceptable_payment_type VARCHAR(100),    -- e.g. Crédit bancaire, cash
    max_year_built INT,                      -- Max acceptable year of construction
    purchase_urgency TEXT,                   -- e.g. Immédiat, dans 1 mois, 3 mois, plus
    comments TEXT,                           -- Internal notes or remarks about the client
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    
    FOREIGN KEY(user_id) REFERENCES users(id)
);
