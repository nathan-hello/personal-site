CREATE TABLE IF NOT EXISTS Images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image TEXT NOT NULL,
    size INTEGER NOT NULL,
    ext TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL,
    author TEXT NOT NULL,
    text TEXT NOT NULL,
    html TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    image_id INTEGER,
    FOREIGN KEY (image_id) REFERENCES Images(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CommentReplies (
    comment_id INTEGER,
    reply_comment_id INTEGER,
    PRIMARY KEY (comment_id, reply_comment_id),
    FOREIGN KEY (comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_comment_id) REFERENCES Comments(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS profiles (
    id TEXT PRIMARY KEY NOT NULL,
    username TEXT UNIQUE NOT NULL,
    color TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chatrooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    creator TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
     
);

CREATE TABLE IF NOT EXISTS chatroom_members (
    chatroom_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (chatroom_id, user_id),
    FOREIGN KEY (chatroom_id) REFERENCES chatrooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id TEXT NOT NULL, -- nullable for anon messages
    message TEXT NOT NULL,
    room_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES chatrooms(id) ON DELETE CASCADE
);
