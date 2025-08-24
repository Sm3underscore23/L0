-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(36) PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR(64) NOT NULL,
    phone VARCHAR(19) NOT NULL,
    zip VARCHAR(9) NOT NULL,
    city VARCHAR(64) NOT NULL,
    address VARCHAR(128) NOT NULL,
    region VARCHAR(64) NOT NULL,
    email VARCHAR(128) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE delivery;
-- +goose StatementEnd
