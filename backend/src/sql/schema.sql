CREATE TABLE users (
    id uuid PRIMARY KEY,
    
    role text NOT NULL,
    name text NOT NULL,
    avt_code text NOT NULL,
    email text UNIQUE,
    
    caro_rating integer NOT NULL DEFAULT 1000,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE caro_matches (
    id uuid PRIMARY KEY,

    x_player_id uuid NOT NULL REFERENCES users(id),
    x_player_rating_before integer NOT NULL,
    x_player_rating_after integer NOT NULL,

    o_player_id uuid NOT NULL REFERENCES users(id),
    o_player_rating_before integer NOT NULL,
    o_player_rating_after integer NOT NULL,

    winner_id uuid REFERENCES users(id),
    final_board jsonb NOT NULL,
    moves jsonb NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_caro_matches_x_player_id
ON caro_matches (x_player_id);

CREATE INDEX idx_caro_matches_o_player_id
ON caro_matches (o_player_id);