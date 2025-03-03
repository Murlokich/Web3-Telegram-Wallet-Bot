BEGIN;

CREATE TABLE IF NOT EXISTS user_master_key (
    user_id BIGINT PRIMARY KEY,
    master_key BYTEA,
    nonce BYTEA
);

CREATE TABLE IF NOT EXISTS user_change_level_key (
    user_id BIGINT PRIMARY KEY,
    change_level_key BYTEA,
    nonce BYTEA
);

COMMIT;