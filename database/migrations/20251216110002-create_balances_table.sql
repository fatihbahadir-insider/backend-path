-- +migrate Up
CREATE TABLE balances (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    amount decimal(15,2) NOT NULL DEFAULT 0,
    last_updated_at timestamp with time zone DEFAULT now(),
    
    CONSTRAINT balances_amount_non_negative CHECK (amount >= 0)
);

-- +migrate Down
DROP TABLE balances;