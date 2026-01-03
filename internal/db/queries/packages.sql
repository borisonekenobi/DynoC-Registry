-- name: CreatePackage :one
INSERT INTO packages (name, description, visibility, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *, (SELECT users.username
              FROM users
              WHERE users.id = $4) AS owner_username;

-- name: GetPackageByID :one
SELECT package.id          AS package_id,
       package.name        AS package_name,
       package.description AS package_description,
       package.visibility  AS package_visibility,
       package.owner_id    AS package_owner_id,
       package.created_at  AS package_created_at,
       package.updated_at  AS package_updated_at
FROM packages package
WHERE id = $1;

-- name: GetPackageByName :one
SELECT package.id          AS package_id,
       package.name        AS package_name,
       package.description AS package_description,
       package.visibility  AS package_visibility,
       package.owner_id    AS package_owner_id,
       package.created_at  AS package_created_at,
       package.updated_at  AS package_updated_at
FROM packages package
WHERE name = $1;

-- name: FindPackages :many
SELECT package.id          AS package_id,
       package.name        AS package_name,
       package.description AS package_description,
       package.visibility  AS package_visibility,
       package.created_at  AS package_created_at,
       package.updated_at  AS package_updated_at,
       usr.username        AS package_owner_username
FROM packages package
JOIN users    usr ON package.owner_id = usr.id
WHERE name ILIKE '%' || $1 || '%'
ORDER BY package.created_at DESC
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
