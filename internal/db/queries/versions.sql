-- name: CreatePackageVersion :one
INSERT INTO package_versions (package_id, version, checksum, size_bytes, location)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPackageVersionByID :one
SELECT *
FROM package_versions
WHERE id = $1;

-- name: GetAllPackageVersions :many
SELECT package_version.*
FROM package_versions package_version
         JOIN packages package ON package_version.package_id = package.id
WHERE package.name = $1
ORDER BY package_version.created_at DESC;

-- name: GetPackageVersionsByName :many
SELECT package_version.*
FROM package_versions package_version
         JOIN packages package ON package_version.package_id = package.id
WHERE package.name = $1
ORDER BY package_version.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPackageByVersion :one
SELECT package.*, package_version.*
FROM packages package
         JOIN package_versions package_version ON package_version.package_id = package.id
WHERE package.name = $1 AND package_version.version = $2;

-- name: UpdatePackageVersion :exec
UPDATE package_versions
SET checksum   = COALESCE($2, checksum),
    size_bytes = COALESCE($3, size_bytes),
    location   = COALESCE($4, location),
    updated_at = NOW()
WHERE id = $1;

-- name: DeletePackageVersion :exec
DELETE
FROM package_versions
WHERE id = $1;

-- name: GetLatestPackageVersion :one
SELECT package_version.*
FROM package_versions package_version
         JOIN packages package ON package_version.package_id = package.id
WHERE package.name = $1
ORDER BY package_version.created_at DESC
LIMIT 1;
