-- name: CreatePackage :one
INSERT INTO packages (name, description, visibility, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *, (
    SELECT users.username
    FROM users
    WHERE users.id = $4
) AS owner_username;

-- name: GetPackageByID :one
SELECT *
FROM packages
WHERE id = $1;

-- name: GetPackageByName :one
SELECT *
FROM packages
WHERE name = $1;

-- name: FindPackages :many
SELECT *
FROM packages
WHERE name ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePackage :exec
UPDATE packages
SET name        = COALESCE($2, name),
    description = COALESCE($3, description),
    visibility  = COALESCE($4, visibility),
    updated_at  = NOW()
WHERE id = $1;

-- name: DeletePackage :exec
DELETE
FROM packages
WHERE id = $1;
