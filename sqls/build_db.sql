CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    hash VARCHAR UNIQUE NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    salt VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL,
    seller_id uuid,
    transaction_address VARCHAR NOT NULL,
    withdraw_address VARCHAR,
    fiat_amount NUMERIC DEFAULT 0.00 NOT NULL,
    deposit_amount NUMERIC DEFAULT 0.00 NOT NULL,
    deposit_exp_date TIMESTAMPTZ NOT NULL,
    fees NUMERIC DEFAULT 0.00 NOT NULL,
    exp_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()::timestamptz,
    active BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO users (hash, username, password, salt) VALUES (
    'c296171f39c1d1801864ff11a8467a0b88f96c2592d86c6ad31210324a71aa96', -- cc177039-1f4d-4d0c-8851-1beffd42ceac
    'seller', 
    '$argon2id$v=19$m=65536,t=3,p=2$cQ9W81MEJ3PWgvv0HJ+pkg$y9JnqFd95BGFXwF3EZ2z8xa/iAmIGQPDqgi8MjzAOAA', 
    'cQ9W81MEJ3PWgvv0HJ+pkg'
) ON CONFLICT DO NOTHING;

INSERT INTO users (hash, username, password, salt) VALUES (
    '5cb46f3cacdd4ada014170108fe11b1d5e3c414ebef25ceaf6e7f42c3388d4ae', -- c75ced33-3534-4c74-b236-2324a4f94a16
    'buyer', 
    '$argon2id$v=19$m=65536,t=3,p=2$KpJ2nx+tZOnu5nsOnKYSQQ$AeuJ4/yDm0aJ8zWDFSMjJnr64oijr3vPDN/JLCvjHb4', 
    'KpJ2nx+tZOnu5nsOnKYSQQ'
) ON CONFLICT DO NOTHING;