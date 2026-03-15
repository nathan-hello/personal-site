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

-- table: profiles
-- name: InsertUserProfile :one
INSERT INTO profiles (id, username, color) VALUES (?, ?, ?) RETURNING *;
-- name: SelectUserProfileById :one
SELECT * FROM profiles WHERE id = ?;
-- name: SelectUserProfileByUsername :one
SELECT * FROM profiles WHERE username = ?;
-- name: UpdateUserProfile :one
UPDATE profiles SET username = ?, color = ? WHERE id = ? RETURNING *;
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
SELECT m.*, p.*
FROM messages AS m
JOIN profiles AS p ON p.id = m.author_id
LEFT JOIN chatroom_members AS cm ON cm.chatroom_id = m.room_id
WHERE m.room_id = ?
ORDER BY m.created_at DESC
LIMIT ?;
SELECT * FROM messages WHERE room_id = ? ORDER BY created_at DESC LIMIT ?;
-- name: SelectMessagesByUser :many
SELECT * FROM messages WHERE author_id = ? ORDER BY created_at DESC LIMIT ?;
-- name: InsertMessage :exec
INSERT INTO messages (author_id, message, room_id, created_at) VALUES (?, ?, ?, ?);
-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ?;
-- name: UpdateMessage :one
UPDATE messages SET message = ? WHERE id = ? RETURNING *;


-- table: chatroom_members
-- name: InsertChatroomMember :exec
INSERT OR IGNORE INTO chatroom_members (chatroom_id, user_id) VALUES (?, ?);
-- name: SelectAllMembersByChatroom :many
SELECT profiles.* 
FROM chatroom_members 
JOIN profiles ON chatroom_members.user_id = profiles.id 
WHERE chatroom_members.chatroom_id = ?;
-- name: SelectUsersJoinedChatrooms :many
SELECT chatroom_members.chatroom_id
FROM chatroom_members 
JOIN chatrooms ON chatroom_members.chatroom_id = chatrooms.id 
WHERE chatroom_members.user_id = ?;
-- DeleteChatroomMember :exec
DELETE FROM chatroom_members WHERE chatroom_members.user_id = ? AND chatroom_members.chatroom_id = ?;

