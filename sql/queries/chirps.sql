-- name: CreateChirp :one
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

-- name: GetChirpsByUserId :many
SELECT * FROM chirps
WHERE user_id=$1;


-- name: GetChirpById :one
SELECT * FROM chirps
WHERE id=$1;

-- name: GetChirpByIdAndUserId :one
SELECT * FROM chirps
WHERE id = $1 and user_id = $2;

-- name: DeleteChirpById :exec
DELETE FROM chirps
WHERE id = $1;


-- name: DeleteChirps :exec
DELETE FROM chirps;

