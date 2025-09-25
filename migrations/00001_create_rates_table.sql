-- +goose Up
CREATE TABLE IF NOT EXISTS rates (
                                     id SERIAL PRIMARY KEY,
                                     currency_code VARCHAR(10) NOT NULL,
                                     price DECIMAL(15, 2) NOT NULL,
                                     timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS rates;