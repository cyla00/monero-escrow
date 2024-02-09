CREATE TABLE IF NOT EXISTS users (
    id uuid UNIQUE NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id uuid UNIQUE NOT NULL,
    owner_id uuid UNIQUE NOT NULL,
    seller_id uuid UNIQUE NOT NULL,
    wallet_address VARCHAR NOT NULL,
    balance NUMERIC DEFAULT 0.00 NOT NULL,
    fees NUMERIC DEFAULT 0.00 NOT NULL,
    exp_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id)
);