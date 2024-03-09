
CREATE TABLE tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    jwt_type TEXT NOT NULL,
    jwt TEXT NOT NULL UNIQUE, 
    valid BOOLEAN NOT NULL,
    family TEXT NOT NULL
);

CREATE TABLE users (
    created_at TEXT DEFAULT (datetime('now')) NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE, 
    encrypted_password TEXT NOT NULL,
    password_created_at TEXT NOT NULL, 
    id TEXT PRIMARY KEY DEFAULT (hex(randomblob(16)))
);

CREATE TABLE users_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    token_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (token_id) REFERENCES tokens(id) ON UPDATE CASCADE ON DELETE CASCADE
);
