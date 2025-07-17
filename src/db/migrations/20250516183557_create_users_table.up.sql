
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fullname TEXT NOT NULL,
    role INT NOT NULL,
    root_id INT NOT,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    agency_name TEXT NOT NULL,
    agency_address TEXT NOT NULL,
    agency_logo TEXT,
    wilaya TEXT NOT NULL,
    daira TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);