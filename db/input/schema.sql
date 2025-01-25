CREATE TABLE IF NOT EXISTS Comments (
    id TEXT PRIMARY KEY DEFAULT (hex(randomblob(16))),
    created_at TEXT DEFAULT "anon" NOT NULL,
    text TEXT NOT NULL,
    post TEXT NOT NULL
);
