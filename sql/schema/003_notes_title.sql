-- +goose Up
ALTER TABLE notes 
ADD title TEXT NOT NULL;

-- +goose Down
ALTER TABLE notes
DROP COLUMN title;