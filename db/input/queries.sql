-- table: Comments
-- name: SelectCommentsMany :many
SELECT * FROM Comments WHERE post_id = ?;
-- name: InsertComment :one
INSERT INTO Comments (author,created_at,text,post_id,html) values (?,?,?,?,?) RETURNING *;

-- table: users
-- name: InsertUser :exec
INSERT INTO users (id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color)
VALUES (?, ?, ?, ?, ?, ?, ?);
-- name: SelectUserByEmail :one
SELECT id, email, username, global_chat_color FROM users WHERE email = ?;
-- name: SelectUserByUsername :one
SELECT id, email, username, global_chat_color FROM users WHERE username = ?;
-- name: SelectUserById :one
SELECT id, email, username, global_chat_color FROM users WHERE id = ?;
-- name: SelectUserByEmailWithPassword :one
SELECT * FROM users WHERE email = ?;
-- name: SelectUserByUsernameWithPassword :one
SELECT * FROM users WHERE username = ?;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- table: users_sessions
-- name: InsertSession :exec
INSERT INTO users_sessions (user_id, token, expires_at, ip_address) VALUES (?, ?, ?, ?);
-- name: SelectSessionByToken :one
SELECT id, user_id, token, created_at, expires_at, last_used, ip_address FROM users_sessions WHERE token = ?;
-- name: SelectSessionsByUser :all
SELECT id, token, created_at, expires_at, last_used, ip_address FROM users_sessions WHERE user_id = ?;
-- name: UpdateSessionLastUsed :exec
UPDATE users_sessions SET last_used = datetime('now') WHERE token = ?;
-- name: DeleteSession :exec
DELETE FROM users_sessions WHERE token = ?;
-- name: DeleteSessionsByUser :exec
DELETE FROM users_sessions WHERE user_id = ?;


-- table: chatrooms
-- name: SelectChatrooms :many
SELECT * FROM chatrooms ORDER BY created_at DESC LIMIT ?;
-- name: InsertChatroom :one
INSERT INTO chatrooms (name, creator, created_at) VALUES (?, ?, ?) RETURNING id;
-- name: DeleteChatroom :exec
DELETE FROM chatrooms WHERE id = ?;
-- name: UpdateChatroomName :one
UPDATE chatrooms SET name = ? WHERE id = ? RETURNING *;


-- table: messages
-- name: SelectMessagesByChatroom :many
SELECT messages.*, chatroom_members.chatroom_color FROM messages
LEFT JOIN chatroom_members ON messages.room_id = chatroom_members.chatroom_id
WHERE messages.room_id = ?
ORDER BY messages.created_at DESC
LIMIT ?;
SELECT * FROM messages WHERE room_id = ? ORDER BY created_at DESC LIMIT ?;
-- name: SelectMessagesByUser :many
SELECT * FROM messages WHERE author_id = ? ORDER BY created_at DESC LIMIT ?;
-- name: InsertMessage :exec
INSERT INTO messages (author_id, author_username, message, room_id, created_at) VALUES (?, ?, ?, ?, ?);
-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ?;
-- name: UpdateMessage :one
UPDATE messages SET message = ? WHERE id = ? RETURNING *;


-- table: chatroom_members
-- name: InsertChatroomMember :exec
INSERT OR IGNORE INTO chatroom_members (chatroom_id, user_id, chatroom_color) VALUES (?, ?, ?);
-- name: SelectAllMembersByChatroom :many
SELECT users.id, users.username, chatroom_members.chatroom_color 
FROM chatroom_members 
JOIN users ON chatroom_members.user_id = users.id 
WHERE chatroom_members.chatroom_id = ?;
-- name: SelectUsersJoinedChatrooms :many
SELECT chatroom_members.chatroom_color, chatroom_members.chatroom_id
FROM chatroom_members 
JOIN chatrooms ON chatroom_members.chatroom_id = chatrooms.id 
WHERE chatroom_members.user_id = ?;
-- DeleteChatroomMember :exec
DELETE FROM chatroom_members WHERE chatroom_members.user_id = ? AND chatroom_members.chatroom_id = ?;
-- name: SelectColorFromUserAndRoom :one
SELECT chatroom_color FROM chatroom_members WHERE chatroom_id = ? AND user_id = ?;



