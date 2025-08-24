-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(36) PRIMARY KEY,
    track_number VARCHAR(36) UNIQUE NOT NULL,
    entry VARCHAR(32) NOT NULL,
    locale VARCHAR(16) NOT NULL,
    internal_signature TEXT NOT NULL,
    customer_id VARCHAR(64) NOT NULL,
    delivery_service VARCHAR(64) NOT NULL,
    shardkey VARCHAR(16) NOT NULL,
    sm_id BIGINT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(16) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
