-- name: CreateDependency :exec
INSERT INTO dependencies (version_id, name, constraint_expr)
VALUES ($1, $2, $3);

-- name: GetDependenciesByVersionID :many
SELECT dependency.name            AS dependency_name,
       dependency.constraint_expr AS constraint_expr
FROM dependencies dependency
WHERE version_id = $1;

-- name: DeleteDependenciesByVersionID :exec
DELETE
FROM dependencies
WHERE version_id = $1;
