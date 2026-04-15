-- +goose Up
CREATE TABLE refresh_tokens(
  token TEXT PRIMARY KEY,
  user_id UUID NOT NULL,
  revoked_at TIMESTAMP,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;
