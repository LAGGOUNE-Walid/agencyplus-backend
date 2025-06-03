CREATE TABLE IF NOT EXISTS sms_queues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,                  
    title VARCHAR(255) NOT NULL,                        
    content TEXT NOT NULL,                     
    from_number VARCHAR(20) NOT NULL,                   
    status VARCHAR(20) DEFAULT 'pending',      -- pending, processing, sent, failed
    total_recipients INTEGER NOT NULL ,        
    sent_count INTEGER DEFAULT 0,              
    failed_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,

    FOREIGN KEY(user_id) REFERENCES users(id)
);
