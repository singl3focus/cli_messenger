-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.tbl_user (
    id serial primary key,
    "name" not null varchar,
    email not null unique varchar,
    password_hash not null varchar,
    "role" not null int,
    created_at not null varchar,
    updated_at not null varchar,
);

COMMENT ON COLUMN auth.tbl_user."role" IS '0 - User, 1 - Admin'  
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE auth.tbl_user;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
