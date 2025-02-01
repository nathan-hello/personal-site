-- table: Comments
-- name: SelectCommentsMany :many
SELECT * FROM Comments WHERE post_id = ?;
-- name: InsertComment :one
INSERT INTO Comments (author,created_at,text,post_id,html) values (?,?,?,?,?) RETURNING *;
