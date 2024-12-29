CREATE TABLE IF NOT EXISTS currencies (
  id SERIAL PRIMARY KEY, -- Unique identifier for each currency
  code CHAR(3) UNIQUE NOT NULL, -- ISO 4217 currency code (e.g., NGN, USD)
  name VARCHAR(50) NOT NULL, -- Full name of the currency (e.g., Naira, Dollar)
  symbol VARCHAR(4) DEFAULT '', -- Symbol for the currency (e.g., ₦, $, €)
  exchange_rate DECIMAL(10, 4) NOT NULL, -- Exchange rate relative to the base currency
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW (),
  base_currency CHAR(3), -- Base currency code for exchange rate reference
  CONSTRAINT fk_base_currency FOREIGN KEY (base_currency) REFERENCES currencies (code) -- Self-referencing FK for base currency
  ON DELETE SET NULL -- Handle cascading deletes appropriately
);

-- default values creations
INSERT INTO
  currencies (code, name, symbol, exchange_rate, base_currency)
VALUES
  ('USD', 'US Dollar', '$', 1.0000, NULL), -- Base currency
  ('NGN', 'Nigerian Naira', '₦', 1644.1400, 'USD'),
  ('GBP', 'British Pound', '£', 0.8154, 'USD'),
  ('EUR', 'Euro', '€', 0.9814, 'USD'),
  ('CAD', 'Canadian Dollar', '$', 1.5145, 'USD');
