CREATE TABLE IF NOT EXISTS Comments (
    id TEXT PRIMARY KEY DEFAULT (hex(randomblob(16))),
    created_at TEXT NOT NULL,
    author TEXT NOT NULL,
    text TEXT NOT NULL,
    post_id TEXT NOT NULL
);
