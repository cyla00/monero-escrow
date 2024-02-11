CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    hash VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    salt VARCHAR NOT NULL,
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

INSERT INTO users (username, password, salt) VALUES (
    'seller', 
    '$argon2id$v=19$m=65536,t=3,p=2$cQ9W81MEJ3PWgvv0HJ+pkg$y9JnqFd95BGFXwF3EZ2z8xa/iAmIGQPDqgi8MjzAOAA', 
    'cQ9W81MEJ3PWgvv0HJ+pkg'
) ON CONFLICT DO NOTHING;

INSERT INTO users (username, password, salt) VALUES (
    'buyer', 
    '$argon2id$v=19$m=65536,t=3,p=2$KpJ2nx+tZOnu5nsOnKYSQQ$AeuJ4/yDm0aJ8zWDFSMjJnr64oijr3vPDN/JLCvjHb4', 
    'KpJ2nx+tZOnu5nsOnKYSQQ'
) ON CONFLICT DO NOTHING;