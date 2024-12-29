CREATE TABLE IF NOT EXISTS wallets (
  id SERIAL PRIMARY KEY,
  public_id VARCHAR(50) UNIQUE NOT NULL,
  balance DECIMAL(20, 2) DEFAULT 0.0,
  is_frozen BOOLEAN DEFAULT false,
  user_id INT REFERENCES users (id) NOT NULL,
  currency_id INT REFERENCES currencies (id) NOT NULL,
  created_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW (),
    updated_at timestamp(0)
  with
    time zone NOT NULL DEFAULT NOW ()
    -- UNIQUE (user_id, currency_id) same as alter below
);

ALTER TABLE wallets ADD CONSTRAINT unique_user_currency UNIQUE (user_id, currency_id);
