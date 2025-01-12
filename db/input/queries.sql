-- table: Users
-- name: SelectAllUsers :many
SELECT * FROM Users;
-- name: InsertUser :one
INSERT INTO Users ( email, username, encrypted_password, password_created_at)
values (?, ?, ?, ?)
RETURNING id, email, username;
-- name: SelectUserByEmail :one
SELECT * FROM Users WHERE email = ?;
-- name: SelectUserByUsername :one
SELECT * FROM Users WHERE username = ?;
-- name: DeleteUser :exec
DELETE FROM Users WHERE id = ?; 
