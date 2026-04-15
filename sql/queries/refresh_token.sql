-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens(token, user_id, expires_at)
VALUES(
  $1,
  $2,
  $3
);

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshTokenByToken :exec
UPDATE  refresh_tokens
SET revoked_at = $1, updated_at = $2
WHERE token = $3;

-- name: DeleteRefreshTokens :exec
DELETE FROM refresh_tokens;
