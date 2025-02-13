CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS balances (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    money INT NOT NULL CHECK (money > 0) DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_name VARCHAR(100) NOT NULL,
    receiver_name VARCHAR(100) NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_inventory (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (user_id, item_name)
);

CREATE OR REPLACE FUNCTION create_user_with_balance(username TEXT, password TEXT)
RETURNS INT AS $$
DECLARE
user_id INT;
BEGIN
INSERT INTO users (username, password)
VALUES (username, password)
    RETURNING id INTO user_id;

INSERT INTO balances (user_id, money)
VALUES (user_id, 1000);

RETURN user_id;
END;
$$ LANGUAGE plpgsql;