CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description)
VALUES
    ('user', 1, 'Regular user'),
    ('moderator', 50, 'Moderator'),
    ('admin', 100, 'Administrator');
