-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

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

-- name: UpdateUserCaroRating :exec
UPDATE users
SET
    caro_rating = GREATEST(0, caro_rating + $2),
    updated_at = now()
WHERE id = $1;

-- name: GetCaroMatchById :one
SELECT * FROM caro_matches
WHERE id = $1 LIMIT 1;

-- name: InsertCaroMatch :exec
INSERT INTO caro_matches (
    id,
    x_player_id,
    x_player_rating_before,
    x_player_rating_after,
    o_player_id,
    o_player_rating_before,
    o_player_rating_after,
    winner_id,
    final_board,
    moves
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
);