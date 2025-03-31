-- table: Comments
-- name: SelectCommentsMany :many
SELECT * FROM Comments WHERE post_id = ?;
-- name: SelectFromComment :one
SELECT * FROM Comments Where id = ?;
-- name: InsertComment :one
INSERT INTO Comments (author,created_at,text,post_id,html,image_id) values (?,?,?,?,?,?) RETURNING *;
-- name: DeleteCommentById :exec
DELETE FROM Comments WHERE id = ?;

-- table: Images
-- name: DeleteImageById :exec
DELETE FROM Images WHERE id = ?;
-- name: SelectFromImage :one
SELECT * FROM Images WHERE id = ? LIMIT 1;
-- name: InsertIntoImage :one
INSERT INTO Images (image,size,ext) values (?,?,?) RETURNING *;

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


-- table: tokens
-- name: SelectTokenFromId :one
SELECT * FROM tokens WHERE id = ?;
-- name: SelectTokenFromJwtString :one
SELECT * FROM tokens WHERE jwt = ?;
-- name: InsertToken :one
INSERT INTO tokens (jwt_type, jwt, valid, family, expires_at) VALUES (?, ?, ?, ?, ?) RETURNING *;
-- name: UpdateTokenValid :one
UPDATE tokens SET valid = ? WHERE jwt = ? RETURNING id;
-- name: UpdateUserTokensToInvalid :exec
UPDATE tokens SET valid = FALSE WHERE id IN (
        SELECT token_id FROM users_tokens WHERE user_id = ?
    );
-- name: UpdateTokensFamilyInvalid :exec
UPDATE tokens SET valid = FALSE WHERE family = ?;
-- name: DeleteTokensByUserId :exec
DELETE FROM tokens WHERE id IN (
        SELECT token_id FROM users_tokens WHERE user_id = ?
    );

-- table: users_tokens
-- name: SelectUsersTokens :many
SELECT * FROM users_tokens WHERE user_id = ?;
-- name: SelectUserIdFromToken :one
SELECT user_id FROM users_tokens WHERE token_id = ? LIMIT 1;
-- name: InsertUsersTokens :exec
INSERT INTO users_tokens (user_id, token_id) VALUES (?, ?);

