-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at Timestamp not null,
    updated_at Timestamp not null,
    email varchar(255) UNIQUE not null
);

-- +goose Down
DROP TABLE users;