CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS balances (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    coins INT NOT NULL CHECK (coins >= 0) DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_inventory (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE (user_id, item_name)
);

CREATE TABLE IF NOT EXISTS items_price (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    price INT NOT NULL CHECK (price >= 0)
);

INSERT INTO items_price (name, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);


CREATE OR REPLACE FUNCTION create_user_with_balance(username TEXT, password TEXT)
RETURNS INT AS $$
DECLARE
user_id INT;
BEGIN
INSERT INTO users (username, password)
VALUES (username, password)
    RETURNING id INTO user_id;

INSERT INTO balances (user_id, coins)
VALUES (user_id, 1000);

RETURN user_id;
END;
$$ LANGUAGE plpgsql;