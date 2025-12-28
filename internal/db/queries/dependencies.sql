-- name: CreateDependency :exec
INSERT INTO dependencies (version_id, dependency_name, constraint_expr)
VALUES ($1, $2, $3);

-- name: GetDependenciesByVersionID :many
SELECT dependency_name, constraint_expr
FROM dependencies
WHERE version_id = $1;

-- name: DeleteDependenciesByVersionID :exec
DELETE
FROM dependencies
WHERE version_id = $1;
