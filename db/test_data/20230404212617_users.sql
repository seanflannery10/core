-- migrate:up
INSERT INTO users (name, email, password_hash, activated)
VALUES ('test', 'test@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', 'true');

-- migrate:down
DELETE FROM users WHERE id = 1;