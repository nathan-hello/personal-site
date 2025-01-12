CREATE TABLE IF NOT EXISTS Users (
    created_at TEXT DEFAULT (datetime('now')) NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE, 
    encrypted_password TEXT NOT NULL,
    password_created_at TEXT NOT NULL, 
    id TEXT PRIMARY KEY DEFAULT (hex(randomblob(16)))
);

CREATE TABLE Image (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alt TEXT NOT NULL,
    url TEXT NOT NULL,
    size TEXT NOT NULL,
    ext TEXT NOT NULL,
    filename TEXT NOT NULL,
    fullname TEXT NOT NULL
);

CREATE TABLE Blog (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    markdown TEXT NOT NULL,
    author TEXT NOT NULL,
    date TEXT NOT NULL, -- ISO8601 (YYYY-MM-DD)
    description TEXT NOT NULL,
    hidden BOOLEAN,
    overrideHref TEXT,
    overrideLayout BOOLEAN,
    tags TEXT, -- JSON array
    title TEXT NOT NULL,
    url TEXT NOT NULL
);

-- Many to many because one image could be on multiple posts
-- And one post could also have multiple images
CREATE TABLE BlogsImages (
    post_id INTEGER NOT NULL,
    image_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, image_id),
    FOREIGN KEY (post_id) REFERENCES Blogs (id),
    FOREIGN KEY (image_id) REFERENCES Images (id)
);

