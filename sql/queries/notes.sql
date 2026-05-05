-- name: CreateNote :one
INSERT INTO notes (id, created_at, updated_at, title ,body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, 
    $2,
    $3
) RETURNING *;