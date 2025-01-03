-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at Timestamp not null,
    updated_at Timestamp not null,
    body varchar(255)  not null,
    user_id varchar(255) not null,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;