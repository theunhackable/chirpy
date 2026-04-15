-- name: CreateUser :one
INSERT INTO users(id, email, hashed_password)
VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserById :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;
