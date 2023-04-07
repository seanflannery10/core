-- migrate:up
INSERT INTO users (name, email, password_hash, activated)
VALUES ('activated', 'activated@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', true);

INSERT INTO users (name, email, password_hash, activated)
VALUES ('unactivated', 'unactivated@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', false);

-- migrate:down
DELETE FROM users WHERE id = 1;
DELETE FROM users WHERE id = 2;