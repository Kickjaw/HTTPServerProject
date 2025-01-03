-- +goose Up
alter table users
add hashed_password TEXT NOT NULL default 'unset';

-- +goose Down
alter table users
drop COLUMN hashed_password;