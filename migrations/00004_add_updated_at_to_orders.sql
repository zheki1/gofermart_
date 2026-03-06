-- +goose Up
ALTER TABLE orders
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE orders DROP COLUMN IF EXISTS updated_at;