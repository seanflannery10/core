-- migrate:up
CREATE TABLE IF NOT EXISTS tokens
(
    hash    bytea PRIMARY KEY,
    user_id bigint       NOT NULL REFERENCES users ON DELETE CASCADE,
    active  bool         NOT NULL DEFAULT true,
    expiry  timestamp(0) NOT NULL,
    scope   text         NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS tokens;
