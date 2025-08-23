-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments (
    transaction VARCHAR(36) PRIMARY KEY REFERENCES orders(order_uid),
    request_id TEXT,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(16),
    amount BIGINT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(64) NOT NULL,
    delivery_cost BIGINT,
    goods_total BIGINT,
    custom_fee BIGINT
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
-- +goose StatementEnd
