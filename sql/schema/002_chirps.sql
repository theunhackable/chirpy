-- +goose Up
CREATE TABLE chirps(
id UUID PRIMARY KEY,
user_id UUID NOT NULL,
body TEXT NOT NULL,
created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
