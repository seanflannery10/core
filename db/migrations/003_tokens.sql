-- migrate:up
CREATE TABLE IF NOT EXISTS tokens
(
    scope   text         NOT NULL,
    expiry  timestamp(0) NOT NULL,
    hash    bytea PRIMARY KEY,
    user_id bigint       NOT NULL REFERENCES users ON DELETE CASCADE,
    active  bool         NOT NULL DEFAULT true
);

-- migrate:down
DROP TABLE IF EXISTS tokens;
