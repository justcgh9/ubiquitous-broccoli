CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    handle TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);

CREATE INDEX idx_users_email ON users (email);

CREATE TABLE apps (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
