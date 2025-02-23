BEGIN;

CREATE TABLE IF NOT EXISTS user_master_key (
    user_id BIGINT PRIMARY KEY,
    master_key VARCHAR(64),
    salt VARCHAR(12)
);

CREATE TABLE IF NOT EXISTS user_change_level_key (
    user_id BIGINT PRIMARY KEY,
    change_level_key VARCHAR(64),
    salt VARCHAR(12)
);

COMMIT;