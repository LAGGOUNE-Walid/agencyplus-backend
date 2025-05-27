CREATE TABLE IF NOT EXISTS building_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id BIGINT NOT NULL,
    building_id BIGINT NOT NULL,
    path VARCHAR(255) NOT NULL,
    mimetype VARCHAR(100),
    size INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (building_id) REFERENCES buildings(id) ON DELETE CASCADE
);
