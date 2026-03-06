-- +goose Up
CREATE TABLE IF NOT EXISTS user_balance (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current NUMERIC(15,2) NOT NULL DEFAULT 0,
    withdrawn NUMERIC(15,2) NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_order_number_idx
ON withdrawals(order_number);

-- +goose Down
DROP INDEX IF EXISTS withdrawals_order_number_idx;
DROP TABLE IF EXISTS user_balance;