CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,        -- session ID (random UUID)
    user_id INTEGER NOT NULL,   -- user ID (foreign key)
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
