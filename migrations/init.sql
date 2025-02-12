CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS balances (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    money INT NOT NULL DEFAULT 1000
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