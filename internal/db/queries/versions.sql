-- name: CreateVersion :one
INSERT INTO package_versions (id, package_id, version, checksum, size_bytes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetLatestVersion :one
SELECT *
FROM package_versions
WHERE package_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: ListVersions :many
SELECT version
FROM package_versions
WHERE package_id = $1
ORDER BY created_at DESC;
