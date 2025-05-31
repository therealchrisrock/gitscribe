-- Revert users table back to numeric ID + firebase_uid structure

-- Create the old table structure
CREATE TABLE users_old (
    id SERIAL PRIMARY KEY,
    firebase_uid VARCHAR(128) UNIQUE,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Copy data back, using the string id as firebase_uid
INSERT INTO users_old (firebase_uid, name, email, created_at, updated_at, deleted_at)
SELECT 
    id as firebase_uid,
    name,
    email,
    created_at,
    updated_at,
    deleted_at
FROM users;

-- Drop the current table
DROP TABLE users;

-- Rename the old table back
ALTER TABLE users_old RENAME TO users;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users(firebase_uid);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add comment
COMMENT ON TABLE users IS 'Stores user information with Firebase Auth integration'; 