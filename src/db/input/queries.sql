-- table: users
-- name: SelectAllUsers :many
SELECT * FROM users;
-- name: InsertUser :one
INSERT INTO users ( email, username, encrypted_password, password_created_at)
values (?, ?, ?, ?)
RETURNING id, email, username;
-- name: SelectUserByEmail :one
SELECT * FROM users WHERE email = ?;
-- name: SelectUserByUsername :one
SELECT * FROM users WHERE username = ?;
-- name: SelectEmailOrUsernameAlreadyExists :one
SELECT email FROM users WHERE users.email = ? OR users.username = ?;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?; 

-- table: tokens
-- name: SelectTokenFromId :one
SELECT * FROM tokens WHERE id = ?;
-- name: SelectTokenFromJwtString :one
SELECT * FROM tokens WHERE jwt = ?;
-- name: InsertToken :one
INSERT INTO tokens (jwt_type, jwt, valid, family) VALUES (?, ?, ?, ?) RETURNING id;
-- name: UpdateTokenValid :one
UPDATE tokens SET valid = ? WHERE jwt = ? RETURNING id;
-- name: UpdateUserTokensToInvalid :exec
UPDATE tokens
SET valid = FALSE
WHERE id IN (
    SELECT token_id
    FROM users_tokens
    WHERE user_id = ?
);
-- name: UpdateTokensFamilyInvalid :exec
UPDATE tokens SET valid = FALSE WHERE family = ?;
-- name: DeleteTokensByUserId :exec
DELETE FROM tokens
WHERE tokens.id IN (
        SELECT token_id FROM users_tokens WHERE users_tokens.user_id = ?
    );

-- table: users_tokens
-- name: SelectUsersTokens :many
SELECT * FROM users_tokens WHERE user_id = ?;
-- name: SelectUserIdFromToken :one
SELECT user_id FROM users_tokens WHERE token_id = ? LIMIT 1;
-- name: InsertUsersTokens :exec
INSERT INTO users_tokens (user_id, token_id) VALUES (?, ?); 
