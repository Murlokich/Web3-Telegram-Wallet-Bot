BEGIN;

CREATE TABLE IF NOT EXISTS user_wallet (
    user_id BIGINT PRIMARY KEY,
    master_key BYTEA,
    master_nonce BYTEA,
    change_level_key BYTEA,
    clk_nonce BYTEA,
    current_address_index INT DEFAULT 0,
    last_address_index INT DEFAULT 0
);

COMMIT;