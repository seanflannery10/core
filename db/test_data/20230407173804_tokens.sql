-- migrate:up
INSERT INTO tokens (scope, expiry, hash, user_id, active)
VALUES ('refresh', '4000-01-01T00:00:00Z', '\x0DE8DF4B9463122F58BC9430BB50095E938A33B627D2B2DC91731EE2C399A812', 3, true);

-- migrate:down
