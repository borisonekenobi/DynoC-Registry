-- name: CreatePackage :one
INSERT INTO packages (id, name, description, visibility, owner_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPackageByName :one
SELECT *
FROM packages
WHERE name = $1;
