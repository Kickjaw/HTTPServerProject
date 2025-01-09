-- name: WriteChirpToDB :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: RetrieveByIDChirp :one
SELECT * FROM chirps 
where id = $1;

-- name: RetrieveByAuthor :many
SELECT * FROM chirps 
where user_id = $1
ORDER BY created_at ASC;

-- name: RetrieveChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: DeleteChirp :exec
DELETE From chirps *
where id = $1;