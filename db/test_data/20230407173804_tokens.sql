-- migrate:up
INSERT INTO tokens (scope, expiry, hash, user_id, active)
VALUES ('refresh', '2023-04-07T18:38:32Z', '\x0DE8DF4B9463122F58BC9430BB50095E938A33B627D2B2DC91731EE2C399A812', 1, true);

-- migrate:down
