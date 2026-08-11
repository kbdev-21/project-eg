CREATE TYPE user_role AS ENUM('ADMIN', 'USER', 'GUEST');
CREATE TYPE user_avt_code AS ENUM('BUNNY', 'KITTEN', 'GRIZZLE', 'HAMSTER', 'MONKEY');

CREATE TABLE users (
    id uuid PRIMARY KEY,
    role user_role NOT NULL,
    name text NOT NULL,
    avt_code user_avt_code NOT NULL,
    email text UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);