CREATE TABLE IF NOT EXISTS Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL,
    author TEXT NOT NULL,
    text TEXT NOT NULL,
    html TEXT NOT NULL,
    post_id INTEGER NOT NULL
);
