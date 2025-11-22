-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id SERIAL,
    user_id INT NOT NULL,
    event TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    mail TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;

-- +goose StatementEnd