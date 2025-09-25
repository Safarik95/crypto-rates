-- +goose Up
CREATE INDEX IF NOT EXISTS idx_rates_currency ON rates(currency_code);
CREATE INDEX IF NOT EXISTS idx_rates_timestamp ON rates(timestamp);

-- +goose Down
DROP INDEX IF EXISTS idx_rates_currency;
DROP INDEX IF EXISTS idx_rates_timestamp;