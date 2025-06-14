-- table: Comments
-- name: SelectCommentsMany :many
SELECT * FROM Comments WHERE post_id = ?;
-- name: SelectFromComment :one
SELECT * FROM Comments Where id = ?;
-- name: InsertComment :one
INSERT INTO Comments (author,created_at,text,post_id,html,image_id) VALUES (?,?,?,?,?,?) RETURNING *;
-- name: DeleteCommentById :exec
DELETE FROM Comments WHERE id = ?;

-- table: Images
-- name: DeleteImageById :exec
DELETE FROM Images WHERE id = ?;
-- name: SelectFromImage :one
SELECT * FROM Images WHERE id = ? LIMIT 1;
-- name: InsertIntoImage :one
INSERT INTO Images (image,size,ext) VALUES (?,?,?) RETURNING *;

-- table: CommentReplies
-- name: InsertReply :exec
INSERT INTO CommentReplies (comment_id, reply_comment_id) VALUES (?,?);
-- name: SelectAllReplies :many
SELECT c.* FROM Comments c JOIN CommentReplies cr ON c.id = cr.reply_comment_id WHERE cr.comment_id = ? ORDER BY cr.reply_comment_id ASC;

-- table: auth_users
-- name: InsertAuthUser :exec
INSERT INTO auth_users (id, email, password_salt, encrypted_password, password_created_at)
VALUES (?, ?, ?, ?, ?);
-- name: SelectAuthUserByEmail :one
SELECT * FROM auth_users WHERE email = ?;
-- name: SelectAuthUserById :one
SELECT * FROM auth_users WHERE id = ?;
-- name: DeleteAuthUser :exec
DELETE FROM auth_users WHERE id = ?;
-- name: UpdateAuthUserPassword :exec
UPDATE auth_users SET encrypted_password = ?, password_salt = ? WHERE id = ?;

-- table: profiles
-- name: InsertUserProfile :exec
INSERT INTO profiles (id, username, global_chat_color)
VALUES (?, ?, ?);
-- name: SelectUserProfileById :one
SELECT * FROM profiles WHERE id = ?;
-- name: SelectUserProfileByUsername :one
SELECT * FROM profiles WHERE username = ?;
-- name: UpdateUserProfile :one
UPDATE profiles 
SET username = ?, global_chat_color = ? 
WHERE id = ? 
RETURNING *;

-- name: DeleteUserProfile :exec
DELETE FROM profiles WHERE id = ?;

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
SELECT profiles.id, profiles.username, chatroom_members.chatroom_color 
FROM chatroom_members 
JOIN profiles ON chatroom_members.user_id = profiles.id 
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

-- name: SelectUserActiveTokens :many
SELECT t.* FROM tokens t
JOIN users_tokens ut ON t.id = ut.token_id
WHERE ut.user_id = ? AND t.valid = true
ORDER BY t.expires_at DESC;

-- name: UpdateTokenInvalid :exec
UPDATE tokens SET valid = false WHERE id = ?;

-- name: UpdateAllUserTokensInvalid :exec
UPDATE tokens SET valid = false 
WHERE id IN (
    SELECT token_id FROM users_tokens WHERE user_id = ?
);

-- name: UpdateUserTokensInvalidExceptFamily :exec
UPDATE tokens SET valid = false 
WHERE id IN (
    SELECT token_id FROM users_tokens WHERE user_id = ?
) AND family != ?;

-- name: UpdateTokenExpiry :exec
UPDATE tokens SET expires_at = ? WHERE id = ?;

-- name: SelectTokenByID :one
SELECT * FROM tokens WHERE id = ?;

