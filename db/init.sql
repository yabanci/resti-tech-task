CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  name TEXT,
  balance DECIMAL
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  value DECIMAL,
  account_id INT REFERENCES accounts(id),
  group TEXT,
  account2_id INT,
  date TIMESTAMP
);

CREATE INDEX idx_account_id ON transactions (account_id);
CREATE INDEX idx_account2_id ON transactions (account2_id);
