-- +migrate Up
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username varchar(50) NOT NULL,
    email varchar(100) NOT NULL,
    password_hash varchar(255) NOT NULL,
    role_id smallint NOT NULL DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    
    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_role_check CHECK (role_id BETWEEN 1 AND 3)
);

CREATE INDEX idx_users_role ON users(role_id);

-- +migrate Down
DROP TABLE users;