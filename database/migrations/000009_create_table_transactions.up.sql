CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('SUCCESS', 'PENDING', 'FAILED')),
    amount NUMERIC(10, 2) NOT NULL,
    transaction_date TIMESTAMP NOT NULL
);

CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date);

