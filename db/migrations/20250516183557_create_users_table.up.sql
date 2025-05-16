
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fullname TEXT NOT NULL,
    role INT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    agency_name TEXT NOT NULL,
    agency_address TEXT NOT NULL,
    agency_logo TEXT,
    wilaya TEXT NOT NULL,
    daira TEXT NOT NULL,
    password NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);