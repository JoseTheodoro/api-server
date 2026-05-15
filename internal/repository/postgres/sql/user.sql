-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
uuid, first_name, last_name, email, genre, date_birth, password
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING *;

-- name: UpdateUser :one
UPDATE users
set first_name = $2, updated_at = $3, last_name = $4, email = $5, password = $6, date_birth = $7, genre = $8
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;