-- migrate:up
CREATE TABLE IF NOT EXISTS messages
(
    id         bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    message    varchar(512)                NOT NULL,
    user_id    bigint                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    version    integer                     NOT NULL DEFAULT 1
);

-- migrate:down
DROP TABLE IF EXISTS messages;