-- name: CreateChip :one
INSERT INTO chirps(id, user_id, body)
VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at
ASC;

-- name: GetChirpById :one
SELECT * FROM chirps
WHERE id=$1;

-- name: DeleteChirps :exec
DELETE FROM chirps;

