CREATE TABLE IF NOT EXISTS device_info_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    device_type TEXT NOT NULL,
    browser TEXT NOT NULL,
    browser_version TEXT NOT NULL,
    os TEXT NOT NULL,
    os_version TEXT NOT NULL,
    app_info TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);