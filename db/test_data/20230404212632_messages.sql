-- migrate:up
INSERT INTO messages (id, message, user_id)
VALUES (1, 'First!', 1);

INSERT INTO messages (id, message, user_id)
VALUES (2, 'Testing', 1);

INSERT INTO messages (id, message, user_id)
VALUES (3, 'Testing!', 1);

INSERT INTO messages (id, message, user_id)
VALUES (4, '4th', 1);

-- migrate:down
DELETE FROM messages WHERE id = 1;
DELETE FROM messages WHERE id = 2;
DELETE FROM messages WHERE id = 3;
DELETE FROM messages WHERE id = 4;