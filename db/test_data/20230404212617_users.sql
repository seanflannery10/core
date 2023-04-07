-- migrate:up
INSERT INTO users (name, email, password_hash, activated)
VALUES ('messages', 'messages@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', true);

INSERT INTO users (name, email, password_hash, activated)
VALUES ('unactivated', 'unactivated@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', false);

INSERT INTO users (name, email, password_hash, activated)
VALUES ('activated', 'activated@test.com', '$2a$13$JHR5woNGzCO6MMhChSgs7OtU/vCADtSj/xb3kBT.fDmFVhuFOgISC', true);

-- migrate:down
