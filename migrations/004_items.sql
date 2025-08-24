-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    track_number VARCHAR(36) REFERENCES orders(track_number),
    chrt_id BIGINT NOT NULL,
    price BIGINT NOT NULL,
    rid TEXT NOT NULL,
    name VARCHAR(64) NOT NULL,
    sale BIGINT NOT NULL,
    size TEXT NOT NULL,
    total_price BIGINT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(64) NOT NULL,
    status BIGINT NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
