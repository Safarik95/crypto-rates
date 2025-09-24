CREATE TABLE IF NOT EXISTS rates (
    id SERIAL PRIMARY KEY,
    currency_code VARCHAR(10) NOT NULL,
    price DECIMAL(15,2) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rates_currency ON rates(currency_code);
CREATE INDEX idx_rates_timestamp ON rates(timestamp);