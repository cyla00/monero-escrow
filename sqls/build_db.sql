CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    transaction_url VARCHAR NOT NULL,
    transaction_address VARCHAR NOT NULL,
    owner_address VARCHAR,
    withdraw_address VARCHAR,
    fiat_amount NUMERIC DEFAULT 0.00 NOT NULL,
    deposit_amount NUMERIC DEFAULT 0.00 NOT NULL,
    fees NUMERIC DEFAULT 0.00 NOT NULL,
    exp_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()::timestamp,
    active BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);