-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: InsertUser :exec
INSERT INTO users (
    id,
    role,
    name,
    avt_code,
    email
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: UpdateUser :exec
UPDATE users
SET
    role = $2,
    name = $3,
    avt_code = $4,
    email = $5,
    updated_at = now()
WHERE id = $1;