-- migrate:up
INSERT INTO users (id, name, email, password_hash, activated)
VALUES (1, 'test', 'test@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', 'true');

-- migrate:down
DELETE FROM users WHERE id = 1;