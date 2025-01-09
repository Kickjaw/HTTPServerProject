-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, 
    $2
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users *;

-- name: GetUserByEmail :one
select * from users
where email = $1;

-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at,revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: FindRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

-- name: RevokeRefreshToke :exec
update refresh_tokens
set revoked_at = NOW(), updated_at = Now()
where token = $1;

-- name: UpdateEmailAndPassword :one
update users
set email = $2, hashed_password = $3, updated_at = NOW()
where id = $1
RETURNING *;

-- name: UpgradeUser :exec
update users
set is_chirpy_red = 'TRUE', updated_at = NOW()
where id = $1;
