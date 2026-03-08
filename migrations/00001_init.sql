-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE orders (
    number TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    status TEXT NOT NULL,
    accrual DOUBLE PRECISION DEFAULT 0,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    order_number TEXT NOT NULL,
    sum DOUBLE PRECISION NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_withdrawals_user_id ON withdrawals(user_id);


-- +goose Down
DROP INDEX IF EXISTS idx_orders_user_id;
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;