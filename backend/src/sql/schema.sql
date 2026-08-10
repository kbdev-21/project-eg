CREATE TABLE users (
    id uuid PRIMARY KEY,
    role text NOT NULL,
    name text NOT NULL,
    avt_code text NOT NULL,
    email text UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);