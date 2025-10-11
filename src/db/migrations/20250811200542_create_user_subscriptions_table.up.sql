CREATE TABLE IF NOT EXISTS user_subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    
    status TEXT DEFAULT 'trial' CHECK (status IN ('active', 'cancelled', 'expired', 'trial')),
    
    -- Key dates (stored as TEXT in ISO format: YYYY-MM-DD)
    current_period_start DATETIME NOT NULL,
    current_period_end DATETIME NOT NULL,
    next_billing_date DATETIME,
    trial_start DATETIME,
    trial_end DATETIME,
    
    -- Pricing
    amount REAL NOT NULL,    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
    
);
