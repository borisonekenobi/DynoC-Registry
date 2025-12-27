-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetSecurityInfo :one
SELECT id, password
FROM users
WHERE username = $1
   OR email = $1;

-- name: UpdateUser :exec
UPDATE users
SET username   = COALESCE($2, username),
    email      = COALESCE($3, email),
    password   = COALESCE($4, password),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;
