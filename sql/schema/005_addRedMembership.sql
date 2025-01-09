-- +goose Up
alter table users
add is_chirpy_red BOOLEAN NOT NULL default 'FALSE';

-- +goose Down
alter table users
drop COLUMN is_chirpy_red;