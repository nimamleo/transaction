-- +migrate Up
CREATE TABLE IF NOT EXISTS system_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ledger_id VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL UNIQUE,
    amount BIGINT NOT NULL DEFAULT 1000000000,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_system_accounts_currency ON system_accounts(currency);

-- +migrate Down
DROP INDEX IF EXISTS idx_system_accounts_currency;
DROP TABLE IF EXISTS system_accounts;
