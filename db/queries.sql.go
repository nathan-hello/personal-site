// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"
	"time"
)

const deleteChatroom = `-- name: DeleteChatroom :exec
DELETE FROM chatrooms WHERE id = ?
`

// DeleteChatroom
//
//	DELETE FROM chatrooms WHERE id = ?
func (q *Queries) DeleteChatroom(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteChatroom, id)
	return err
}

const deleteCommentById = `-- name: DeleteCommentById :exec
DELETE FROM Comments WHERE id = ?
`

// DeleteCommentById
//
//	DELETE FROM Comments WHERE id = ?
func (q *Queries) DeleteCommentById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCommentById, id)
	return err
}

const deleteImageById = `-- name: DeleteImageById :exec
DELETE FROM Images WHERE id = ?
`

// table: Images
//
//	DELETE FROM Images WHERE id = ?
func (q *Queries) DeleteImageById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteImageById, id)
	return err
}

const deleteMessage = `-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ?
`

// DeleteMessage
//
//	DELETE FROM messages WHERE id = ?
func (q *Queries) DeleteMessage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMessage, id)
	return err
}

const deleteTokensByUserId = `-- name: DeleteTokensByUserId :exec
DELETE FROM tokens WHERE id IN (
        SELECT token_id FROM users_tokens WHERE user_id = ?
    )
`

// DeleteTokensByUserId
//
//	DELETE FROM tokens WHERE id IN (
//	        SELECT token_id FROM users_tokens WHERE user_id = ?
//	    )
func (q *Queries) DeleteTokensByUserId(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, deleteTokensByUserId, userID)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?
`

// DeleteUser
//
//	DELETE FROM users WHERE id = ?
func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const insertChatroom = `-- name: InsertChatroom :one
INSERT INTO chatrooms (name, creator, created_at) VALUES (?, ?, ?) RETURNING id
`

type InsertChatroomParams struct {
	Name      string
	Creator   string
	CreatedAt time.Time
}

// InsertChatroom
//
//	INSERT INTO chatrooms (name, creator, created_at) VALUES (?, ?, ?) RETURNING id
func (q *Queries) InsertChatroom(ctx context.Context, arg InsertChatroomParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertChatroom, arg.Name, arg.Creator, arg.CreatedAt)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertChatroomMember = `-- name: InsertChatroomMember :exec
INSERT OR IGNORE INTO chatroom_members (chatroom_id, user_id, chatroom_color) VALUES (?, ?, ?)
`

type InsertChatroomMemberParams struct {
	ChatroomID    int64
	UserID        string
	ChatroomColor string
}

// table: chatroom_members
//
//	INSERT OR IGNORE INTO chatroom_members (chatroom_id, user_id, chatroom_color) VALUES (?, ?, ?)
func (q *Queries) InsertChatroomMember(ctx context.Context, arg InsertChatroomMemberParams) error {
	_, err := q.db.ExecContext(ctx, insertChatroomMember, arg.ChatroomID, arg.UserID, arg.ChatroomColor)
	return err
}

const insertComment = `-- name: InsertComment :one
INSERT INTO Comments (author,created_at,text,post_id,html,image_id) VALUES (?,?,?,?,?,?) RETURNING id, created_at, author, text, html, post_id, image_id
`

type InsertCommentParams struct {
	Author    string
	CreatedAt string
	Text      string
	PostID    int64
	Html      string
	ImageID   *int64
}

// InsertComment
//
//	INSERT INTO Comments (author,created_at,text,post_id,html,image_id) VALUES (?,?,?,?,?,?) RETURNING id, created_at, author, text, html, post_id, image_id
func (q *Queries) InsertComment(ctx context.Context, arg InsertCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, insertComment,
		arg.Author,
		arg.CreatedAt,
		arg.Text,
		arg.PostID,
		arg.Html,
		arg.ImageID,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Author,
		&i.Text,
		&i.Html,
		&i.PostID,
		&i.ImageID,
	)
	return i, err
}

const insertIntoImage = `-- name: InsertIntoImage :one
INSERT INTO Images (image,size,ext) VALUES (?,?,?) RETURNING id, image, size, ext
`

type InsertIntoImageParams struct {
	Image string
	Size  int64
	Ext   string
}

// InsertIntoImage
//
//	INSERT INTO Images (image,size,ext) VALUES (?,?,?) RETURNING id, image, size, ext
func (q *Queries) InsertIntoImage(ctx context.Context, arg InsertIntoImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, insertIntoImage, arg.Image, arg.Size, arg.Ext)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Image,
		&i.Size,
		&i.Ext,
	)
	return i, err
}

const insertMessage = `-- name: InsertMessage :exec
INSERT INTO messages (author_id, author_username, message, room_id, created_at) VALUES (?, ?, ?, ?, ?)
`

type InsertMessageParams struct {
	AuthorID       *string
	AuthorUsername string
	Message        string
	RoomID         int64
	CreatedAt      time.Time
}

// InsertMessage
//
//	INSERT INTO messages (author_id, author_username, message, room_id, created_at) VALUES (?, ?, ?, ?, ?)
func (q *Queries) InsertMessage(ctx context.Context, arg InsertMessageParams) error {
	_, err := q.db.ExecContext(ctx, insertMessage,
		arg.AuthorID,
		arg.AuthorUsername,
		arg.Message,
		arg.RoomID,
		arg.CreatedAt,
	)
	return err
}

const insertReply = `-- name: InsertReply :exec
INSERT INTO CommentReplies (comment_id, reply_comment_id) VALUES (?,?)
`

type InsertReplyParams struct {
	CommentID      *int64
	ReplyCommentID *int64
}

// table: CommentReplies
//
//	INSERT INTO CommentReplies (comment_id, reply_comment_id) VALUES (?,?)
func (q *Queries) InsertReply(ctx context.Context, arg InsertReplyParams) error {
	_, err := q.db.ExecContext(ctx, insertReply, arg.CommentID, arg.ReplyCommentID)
	return err
}

const insertToken = `-- name: InsertToken :one
INSERT INTO tokens (jwt_type, jwt, valid, family, expires_at) VALUES (?, ?, ?, ?, ?) RETURNING id, jwt_type, jwt, valid, family, expires_at
`

type InsertTokenParams struct {
	JwtType   string
	Jwt       string
	Valid     bool
	Family    string
	ExpiresAt int64
}

// InsertToken
//
//	INSERT INTO tokens (jwt_type, jwt, valid, family, expires_at) VALUES (?, ?, ?, ?, ?) RETURNING id, jwt_type, jwt, valid, family, expires_at
func (q *Queries) InsertToken(ctx context.Context, arg InsertTokenParams) (Token, error) {
	row := q.db.QueryRowContext(ctx, insertToken,
		arg.JwtType,
		arg.Jwt,
		arg.Valid,
		arg.Family,
		arg.ExpiresAt,
	)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.JwtType,
		&i.Jwt,
		&i.Valid,
		&i.Family,
		&i.ExpiresAt,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :exec
INSERT INTO users (id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

type InsertUserParams struct {
	ID                string
	Email             string
	Username          string
	PasswordSalt      string
	EncryptedPassword string
	PasswordCreatedAt time.Time
	GlobalChatColor   string
}

// table: users
//
//	INSERT INTO users (id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color)
//	VALUES (?, ?, ?, ?, ?, ?, ?)
func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) error {
	_, err := q.db.ExecContext(ctx, insertUser,
		arg.ID,
		arg.Email,
		arg.Username,
		arg.PasswordSalt,
		arg.EncryptedPassword,
		arg.PasswordCreatedAt,
		arg.GlobalChatColor,
	)
	return err
}

const insertUsersTokens = `-- name: InsertUsersTokens :exec
INSERT INTO users_tokens (user_id, token_id) VALUES (?, ?)
`

type InsertUsersTokensParams struct {
	UserID  string
	TokenID int64
}

// InsertUsersTokens
//
//	INSERT INTO users_tokens (user_id, token_id) VALUES (?, ?)
func (q *Queries) InsertUsersTokens(ctx context.Context, arg InsertUsersTokensParams) error {
	_, err := q.db.ExecContext(ctx, insertUsersTokens, arg.UserID, arg.TokenID)
	return err
}

const selectAllMembersByChatroom = `-- name: SelectAllMembersByChatroom :many
SELECT users.id, users.username, chatroom_members.chatroom_color 
FROM chatroom_members 
JOIN users ON chatroom_members.user_id = users.id 
WHERE chatroom_members.chatroom_id = ?
`

type SelectAllMembersByChatroomRow struct {
	ID            string
	Username      string
	ChatroomColor string
}

// SelectAllMembersByChatroom
//
//	SELECT users.id, users.username, chatroom_members.chatroom_color
//	FROM chatroom_members
//	JOIN users ON chatroom_members.user_id = users.id
//	WHERE chatroom_members.chatroom_id = ?
func (q *Queries) SelectAllMembersByChatroom(ctx context.Context, chatroomID int64) ([]SelectAllMembersByChatroomRow, error) {
	rows, err := q.db.QueryContext(ctx, selectAllMembersByChatroom, chatroomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectAllMembersByChatroomRow
	for rows.Next() {
		var i SelectAllMembersByChatroomRow
		if err := rows.Scan(&i.ID, &i.Username, &i.ChatroomColor); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectAllReplies = `-- name: SelectAllReplies :many
SELECT c.id, c.created_at, c.author, c.text, c.html, c.post_id, c.image_id FROM Comments c JOIN CommentReplies cr ON c.id = cr.reply_comment_id WHERE cr.comment_id = ? ORDER BY cr.reply_comment_id ASC
`

// SelectAllReplies
//
//	SELECT c.id, c.created_at, c.author, c.text, c.html, c.post_id, c.image_id FROM Comments c JOIN CommentReplies cr ON c.id = cr.reply_comment_id WHERE cr.comment_id = ? ORDER BY cr.reply_comment_id ASC
func (q *Queries) SelectAllReplies(ctx context.Context, commentID *int64) ([]Comment, error) {
	rows, err := q.db.QueryContext(ctx, selectAllReplies, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Author,
			&i.Text,
			&i.Html,
			&i.PostID,
			&i.ImageID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectChatrooms = `-- name: SelectChatrooms :many
SELECT id, name, creator, created_at FROM chatrooms ORDER BY created_at DESC LIMIT ?
`

// table: chatrooms
//
//	SELECT id, name, creator, created_at FROM chatrooms ORDER BY created_at DESC LIMIT ?
func (q *Queries) SelectChatrooms(ctx context.Context, limit int64) ([]Chatroom, error) {
	rows, err := q.db.QueryContext(ctx, selectChatrooms, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chatroom
	for rows.Next() {
		var i Chatroom
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Creator,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectColorFromUserAndRoom = `-- name: SelectColorFromUserAndRoom :one
SELECT chatroom_color FROM chatroom_members WHERE chatroom_id = ? AND user_id = ?
`

type SelectColorFromUserAndRoomParams struct {
	ChatroomID int64
	UserID     string
}

// SelectColorFromUserAndRoom
//
//	SELECT chatroom_color FROM chatroom_members WHERE chatroom_id = ? AND user_id = ?
func (q *Queries) SelectColorFromUserAndRoom(ctx context.Context, arg SelectColorFromUserAndRoomParams) (string, error) {
	row := q.db.QueryRowContext(ctx, selectColorFromUserAndRoom, arg.ChatroomID, arg.UserID)
	var chatroom_color string
	err := row.Scan(&chatroom_color)
	return chatroom_color, err
}

const selectCommentsMany = `-- name: SelectCommentsMany :many
SELECT id, created_at, author, text, html, post_id, image_id FROM Comments WHERE post_id = ?
`

// table: Comments
//
//	SELECT id, created_at, author, text, html, post_id, image_id FROM Comments WHERE post_id = ?
func (q *Queries) SelectCommentsMany(ctx context.Context, postID int64) ([]Comment, error) {
	rows, err := q.db.QueryContext(ctx, selectCommentsMany, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Author,
			&i.Text,
			&i.Html,
			&i.PostID,
			&i.ImageID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectFromComment = `-- name: SelectFromComment :one
SELECT id, created_at, author, text, html, post_id, image_id FROM Comments Where id = ?
`

// SelectFromComment
//
//	SELECT id, created_at, author, text, html, post_id, image_id FROM Comments Where id = ?
func (q *Queries) SelectFromComment(ctx context.Context, id int64) (Comment, error) {
	row := q.db.QueryRowContext(ctx, selectFromComment, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Author,
		&i.Text,
		&i.Html,
		&i.PostID,
		&i.ImageID,
	)
	return i, err
}

const selectFromImage = `-- name: SelectFromImage :one
SELECT id, image, size, ext FROM Images WHERE id = ? LIMIT 1
`

// SelectFromImage
//
//	SELECT id, image, size, ext FROM Images WHERE id = ? LIMIT 1
func (q *Queries) SelectFromImage(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, selectFromImage, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Image,
		&i.Size,
		&i.Ext,
	)
	return i, err
}

const selectMessagesByChatroom = `-- name: SelectMessagesByChatroom :many
SELECT messages.id, messages.author_id, messages.author_username, messages.message, messages.room_id, messages.created_at, chatroom_members.chatroom_color FROM messages
LEFT JOIN chatroom_members ON messages.room_id = chatroom_members.chatroom_id
WHERE messages.room_id = ?
ORDER BY messages.created_at DESC
LIMIT ?
`

type SelectMessagesByChatroomParams struct {
	RoomID int64
	Limit  int64
}

type SelectMessagesByChatroomRow struct {
	ID             int64
	AuthorID       *string
	AuthorUsername string
	Message        string
	RoomID         int64
	CreatedAt      time.Time
	ChatroomColor  *string
}

// table: messages
//
//	SELECT messages.id, messages.author_id, messages.author_username, messages.message, messages.room_id, messages.created_at, chatroom_members.chatroom_color FROM messages
//	LEFT JOIN chatroom_members ON messages.room_id = chatroom_members.chatroom_id
//	WHERE messages.room_id = ?
//	ORDER BY messages.created_at DESC
//	LIMIT ?
func (q *Queries) SelectMessagesByChatroom(ctx context.Context, arg SelectMessagesByChatroomParams) ([]SelectMessagesByChatroomRow, error) {
	rows, err := q.db.QueryContext(ctx, selectMessagesByChatroom, arg.RoomID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectMessagesByChatroomRow
	for rows.Next() {
		var i SelectMessagesByChatroomRow
		if err := rows.Scan(
			&i.ID,
			&i.AuthorID,
			&i.AuthorUsername,
			&i.Message,
			&i.RoomID,
			&i.CreatedAt,
			&i.ChatroomColor,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectMessagesByUser = `-- name: SelectMessagesByUser :many
SELECT id, author_id, author_username, message, room_id, created_at FROM messages WHERE author_id = ? ORDER BY created_at DESC LIMIT ?
`

type SelectMessagesByUserParams struct {
	AuthorID *string
	Limit    int64
}

// SelectMessagesByUser
//
//	SELECT id, author_id, author_username, message, room_id, created_at FROM messages WHERE author_id = ? ORDER BY created_at DESC LIMIT ?
func (q *Queries) SelectMessagesByUser(ctx context.Context, arg SelectMessagesByUserParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, selectMessagesByUser, arg.AuthorID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.AuthorID,
			&i.AuthorUsername,
			&i.Message,
			&i.RoomID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectTokenFromId = `-- name: SelectTokenFromId :one
SELECT id, jwt_type, jwt, valid, family, expires_at FROM tokens WHERE id = ?
`

// table: tokens
//
//	SELECT id, jwt_type, jwt, valid, family, expires_at FROM tokens WHERE id = ?
func (q *Queries) SelectTokenFromId(ctx context.Context, id int64) (Token, error) {
	row := q.db.QueryRowContext(ctx, selectTokenFromId, id)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.JwtType,
		&i.Jwt,
		&i.Valid,
		&i.Family,
		&i.ExpiresAt,
	)
	return i, err
}

const selectTokenFromJwtString = `-- name: SelectTokenFromJwtString :one
SELECT id, jwt_type, jwt, valid, family, expires_at FROM tokens WHERE jwt = ?
`

// SelectTokenFromJwtString
//
//	SELECT id, jwt_type, jwt, valid, family, expires_at FROM tokens WHERE jwt = ?
func (q *Queries) SelectTokenFromJwtString(ctx context.Context, jwt string) (Token, error) {
	row := q.db.QueryRowContext(ctx, selectTokenFromJwtString, jwt)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.JwtType,
		&i.Jwt,
		&i.Valid,
		&i.Family,
		&i.ExpiresAt,
	)
	return i, err
}

const selectUserByEmail = `-- name: SelectUserByEmail :one
SELECT id, email, username, global_chat_color FROM users WHERE email = ?
`

type SelectUserByEmailRow struct {
	ID              string
	Email           string
	Username        string
	GlobalChatColor string
}

// SelectUserByEmail
//
//	SELECT id, email, username, global_chat_color FROM users WHERE email = ?
func (q *Queries) SelectUserByEmail(ctx context.Context, email string) (SelectUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, selectUserByEmail, email)
	var i SelectUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.GlobalChatColor,
	)
	return i, err
}

const selectUserByEmailWithPassword = `-- name: SelectUserByEmailWithPassword :one
SELECT id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color FROM users WHERE email = ?
`

// SelectUserByEmailWithPassword
//
//	SELECT id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color FROM users WHERE email = ?
func (q *Queries) SelectUserByEmailWithPassword(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserByEmailWithPassword, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.PasswordSalt,
		&i.EncryptedPassword,
		&i.PasswordCreatedAt,
		&i.GlobalChatColor,
	)
	return i, err
}

const selectUserById = `-- name: SelectUserById :one
SELECT id, email, username, global_chat_color FROM users WHERE id = ?
`

type SelectUserByIdRow struct {
	ID              string
	Email           string
	Username        string
	GlobalChatColor string
}

// SelectUserById
//
//	SELECT id, email, username, global_chat_color FROM users WHERE id = ?
func (q *Queries) SelectUserById(ctx context.Context, id string) (SelectUserByIdRow, error) {
	row := q.db.QueryRowContext(ctx, selectUserById, id)
	var i SelectUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.GlobalChatColor,
	)
	return i, err
}

const selectUserByUsername = `-- name: SelectUserByUsername :one
SELECT id, email, username, global_chat_color FROM users WHERE username = ?
`

type SelectUserByUsernameRow struct {
	ID              string
	Email           string
	Username        string
	GlobalChatColor string
}

// SelectUserByUsername
//
//	SELECT id, email, username, global_chat_color FROM users WHERE username = ?
func (q *Queries) SelectUserByUsername(ctx context.Context, username string) (SelectUserByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, selectUserByUsername, username)
	var i SelectUserByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.GlobalChatColor,
	)
	return i, err
}

const selectUserByUsernameWithPassword = `-- name: SelectUserByUsernameWithPassword :one
SELECT id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color FROM users WHERE username = ?
`

// SelectUserByUsernameWithPassword
//
//	SELECT id, email, username, password_salt, encrypted_password, password_created_at, global_chat_color FROM users WHERE username = ?
func (q *Queries) SelectUserByUsernameWithPassword(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserByUsernameWithPassword, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.PasswordSalt,
		&i.EncryptedPassword,
		&i.PasswordCreatedAt,
		&i.GlobalChatColor,
	)
	return i, err
}

const selectUserIdFromToken = `-- name: SelectUserIdFromToken :one
SELECT user_id FROM users_tokens WHERE token_id = ? LIMIT 1
`

// SelectUserIdFromToken
//
//	SELECT user_id FROM users_tokens WHERE token_id = ? LIMIT 1
func (q *Queries) SelectUserIdFromToken(ctx context.Context, tokenID int64) (string, error) {
	row := q.db.QueryRowContext(ctx, selectUserIdFromToken, tokenID)
	var user_id string
	err := row.Scan(&user_id)
	return user_id, err
}

const selectUsersJoinedChatrooms = `-- name: SelectUsersJoinedChatrooms :many
SELECT chatroom_members.chatroom_color, chatroom_members.chatroom_id
FROM chatroom_members 
JOIN chatrooms ON chatroom_members.chatroom_id = chatrooms.id 
WHERE chatroom_members.user_id = ?
`

type SelectUsersJoinedChatroomsRow struct {
	ChatroomColor string
	ChatroomID    int64
}

// SelectUsersJoinedChatrooms
//
//	SELECT chatroom_members.chatroom_color, chatroom_members.chatroom_id
//	FROM chatroom_members
//	JOIN chatrooms ON chatroom_members.chatroom_id = chatrooms.id
//	WHERE chatroom_members.user_id = ?
func (q *Queries) SelectUsersJoinedChatrooms(ctx context.Context, userID string) ([]SelectUsersJoinedChatroomsRow, error) {
	rows, err := q.db.QueryContext(ctx, selectUsersJoinedChatrooms, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectUsersJoinedChatroomsRow
	for rows.Next() {
		var i SelectUsersJoinedChatroomsRow
		if err := rows.Scan(&i.ChatroomColor, &i.ChatroomID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectUsersTokens = `-- name: SelectUsersTokens :many
SELECT user_id, token_id FROM users_tokens WHERE user_id = ?
`

// table: users_tokens
//
//	SELECT user_id, token_id FROM users_tokens WHERE user_id = ?
func (q *Queries) SelectUsersTokens(ctx context.Context, userID string) ([]UsersToken, error) {
	rows, err := q.db.QueryContext(ctx, selectUsersTokens, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersToken
	for rows.Next() {
		var i UsersToken
		if err := rows.Scan(&i.UserID, &i.TokenID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateChatroomName = `-- name: UpdateChatroomName :one
UPDATE chatrooms SET name = ? WHERE id = ? RETURNING id, name, creator, created_at
`

type UpdateChatroomNameParams struct {
	Name string
	ID   int64
}

// UpdateChatroomName
//
//	UPDATE chatrooms SET name = ? WHERE id = ? RETURNING id, name, creator, created_at
func (q *Queries) UpdateChatroomName(ctx context.Context, arg UpdateChatroomNameParams) (Chatroom, error) {
	row := q.db.QueryRowContext(ctx, updateChatroomName, arg.Name, arg.ID)
	var i Chatroom
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Creator,
		&i.CreatedAt,
	)
	return i, err
}

const updateMessage = `-- name: UpdateMessage :one
UPDATE messages SET message = ? WHERE id = ? RETURNING id, author_id, author_username, message, room_id, created_at
`

type UpdateMessageParams struct {
	Message string
	ID      int64
}

// UpdateMessage
//
//	UPDATE messages SET message = ? WHERE id = ? RETURNING id, author_id, author_username, message, room_id, created_at
func (q *Queries) UpdateMessage(ctx context.Context, arg UpdateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, updateMessage, arg.Message, arg.ID)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.AuthorUsername,
		&i.Message,
		&i.RoomID,
		&i.CreatedAt,
	)
	return i, err
}

const updateTokenValid = `-- name: UpdateTokenValid :one
UPDATE tokens SET valid = ? WHERE jwt = ? RETURNING id
`

type UpdateTokenValidParams struct {
	Valid bool
	Jwt   string
}

// UpdateTokenValid
//
//	UPDATE tokens SET valid = ? WHERE jwt = ? RETURNING id
func (q *Queries) UpdateTokenValid(ctx context.Context, arg UpdateTokenValidParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateTokenValid, arg.Valid, arg.Jwt)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateTokensFamilyInvalid = `-- name: UpdateTokensFamilyInvalid :exec
UPDATE tokens SET valid = FALSE WHERE family = ?
`

// UpdateTokensFamilyInvalid
//
//	UPDATE tokens SET valid = FALSE WHERE family = ?
func (q *Queries) UpdateTokensFamilyInvalid(ctx context.Context, family string) error {
	_, err := q.db.ExecContext(ctx, updateTokensFamilyInvalid, family)
	return err
}

const updateUserTokensToInvalid = `-- name: UpdateUserTokensToInvalid :exec
UPDATE tokens SET valid = FALSE WHERE id IN (
        SELECT token_id FROM users_tokens WHERE user_id = ?
    )
`

// UpdateUserTokensToInvalid
//
//	UPDATE tokens SET valid = FALSE WHERE id IN (
//	        SELECT token_id FROM users_tokens WHERE user_id = ?
//	    )
func (q *Queries) UpdateUserTokensToInvalid(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, updateUserTokensToInvalid, userID)
	return err
}
