-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.tbl_user (
    id SERIAL primary key,
    "name" VARCHAR not null,
    email VARCHAR not null unique,
    password_hash VARCHAR not null,
    "role" INT not null,
    created_at TIMESTAMP WITH TIME ZONE not null,
    updated_at TIMESTAMP WITH TIME ZONE not null
);

COMMENT ON COLUMN auth.tbl_user."role" IS '0 - User, 1 - Admin';
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE auth.tbl_user;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
