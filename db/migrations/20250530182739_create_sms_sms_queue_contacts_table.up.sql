CREATE TABLE IF NOT EXISTS sms_queue_contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sms_queue_id INTEGER NOT NULL REFERENCES sms_queue(id),
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    error TEXT,                                -- Optional error message
    sent_at TIMESTAMP                          -- Time this number was sent
);
