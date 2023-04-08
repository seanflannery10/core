-- migrate:up
INSERT INTO tokens (scope, expiry, hash, user_id, active)
VALUES ('refresh', '4000-01-01T00:00:00Z', '\x0DE8DF4B9463122F58BC9430BB50095E938A33B627D2B2DC91731EE2C399A812', 3, true);

INSERT INTO tokens (scope, expiry, hash, user_id, active)
VALUES ('activation', '4000-01-01T00:00:00Z', '\x7013DBD4D6857FA93B6766CC7A6CEEA2B1B7323552FBDC9A83C6A10B2C321097', 4, true);

INSERT INTO tokens (scope, expiry, hash, user_id, active)
VALUES ('password-reset', '4000-01-01T00:00:00Z', '\x1F16D1CEE79A483483F2DBB1214983F52549B8A95C5C56B234727494D481370F', 4, true);
-- migrate:down
