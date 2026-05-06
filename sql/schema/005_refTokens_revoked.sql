-- +goose Up
ALTER TABLE refresh_tokens 
ADD revoked_at TIMESTAMP;

-- +goose Down 
DROP COLUMN revoked_at;