ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS role_id BIGINT REFERENCES roles(id) DEFAULT 1;

UPDATE users
SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
)
WHERE role_id IS NULL;

ALTER TABLE users
    ALTER COLUMN role_id DROP DEFAULT;

ALTER TABLE users
    ALTER COLUMN role_id SET NOT NULL;
