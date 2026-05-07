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

-- name: GetNotesByUserID :many
SELECT * FROM notes
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetNoteByID :one
SELECT * FROM notes
WHERE id = $1;

-- name: UpdateNote :one
UPDATE notes
SET
  title = COALESCE(sqlc.narg('title'), title),
  body = COALESCE(sqlc.narg('body'), body),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1 AND user_id = $2;