CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    hash VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    owner_id uuid UNIQUE NOT NULL,
    seller_id uuid UNIQUE NOT NULL,
    wallet_address VARCHAR NOT NULL,
    balance NUMERIC DEFAULT 0.00 NOT NULL,
    fees NUMERIC DEFAULT 0.00 NOT NULL,
    exp_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    active BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO users ( 
    username, 
    password 
)
VALUES (
    'buyeruser', 
    '$argon2id$v=19$m=16,t=2,p=1$cGFzc3dvcmQ$8vDS3rsezOjrur01dF12EA' -- salt: password
) ON CONFLICT DO NOTHING;

INSERT INTO users ( 
    username, 
    password 
)
VALUES (
    'selleruser', 
    '$argon2id$v=19$m=16,t=2,p=1$cGFzc3dvcmQ$8vDS3rsezOjrur01dF12EA' -- salt: password
) ON CONFLICT DO NOTHING;