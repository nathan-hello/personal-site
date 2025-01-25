-- table: Comments
-- name: InsertComment :one
INSERT INTO Comments (created_at,text,post) values (?,?,?) RETURNING *;
