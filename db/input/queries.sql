-- table: Comments
-- name: SelectCommentsMany :many
SELECT * FROM Comments WHERE post_id = ?;
-- name: InsertComment :one
INSERT INTO Comments (created_at,text,post_id) values (?,?,?) RETURNING *;
