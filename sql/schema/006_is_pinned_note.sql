-- +goose Up
ALTER TABLE notes
ADD is_pinned BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE notes 
DROP COLUMN is_pinned;