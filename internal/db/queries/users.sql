-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;

-- name: GetUserByID :one
SELECT usr.id         AS user_id,
       usr.username   AS user_name,
       usr.email      AS user_email,
       usr.created_at AS user_created_at,
       usr.updated_at AS user_updated_at
FROM users usr
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT usr.id         AS user_id,
       usr.username   AS user_name,
       usr.email      AS user_email,
       usr.created_at AS user_created_at,
       usr.updated_at AS user_updated_at
FROM users usr
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT usr.id         AS user_id,
       usr.username   AS user_name,
       usr.email      AS user_email,
       usr.created_at AS user_created_at,
       usr.updated_at AS user_updated_at
FROM users usr
WHERE email = $1;

-- name: GetSecurityInfo :one
SELECT usr.id       AS user_id,
       usr.password AS user_password
FROM users usr
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
