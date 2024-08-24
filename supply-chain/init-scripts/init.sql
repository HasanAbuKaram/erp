-- Connect to the default database
\c postgres

-- Drop unwanted databases
-- Avoid dropping default databases template0 and template1
-- Drop any other databases you no longer need
-- DROP DATABASE IF EXISTS some_other_database;

-- Optionally, you can also drop roles if not needed
-- DROP ROLE IF EXISTS some_role;

-- Create the items table in the supply_chain database
\c supply_chain

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
    