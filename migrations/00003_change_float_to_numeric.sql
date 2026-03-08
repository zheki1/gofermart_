-- +goose Up
ALTER TABLE orders
ALTER COLUMN accrual TYPE NUMERIC(15,2)
USING ROUND(accrual::numeric, 2);

ALTER TABLE withdrawals
ALTER COLUMN sum TYPE NUMERIC(15,2)
USING ROUND(sum::numeric, 2);

-- +goose Down
ALTER TABLE orders
ALTER COLUMN accrual TYPE DOUBLE PRECISION
USING accrual::double precision;

ALTER TABLE withdrawals
ALTER COLUMN sum TYPE DOUBLE PRECISION
USING sum::double precision;