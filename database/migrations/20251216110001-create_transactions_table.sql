-- +migrate Up
CREATE TABLE transactions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    to_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    amount decimal(15,2) NOT NULL,
    type smallint NOT NULL,
    status smallint NOT NULL DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    
    CONSTRAINT transactions_type_check CHECK (type BETWEEN 1 AND 3),
    CONSTRAINT transactions_status_check CHECK (status BETWEEN 1 AND 4),
    CONSTRAINT transactions_amount_positive CHECK (amount > 0)
);

CREATE INDEX idx_transactions_from_user ON transactions(from_user_id);
CREATE INDEX idx_transactions_to_user ON transactions(to_user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created ON transactions(created_at);

-- +migrate Down
DROP TABLE transactions;