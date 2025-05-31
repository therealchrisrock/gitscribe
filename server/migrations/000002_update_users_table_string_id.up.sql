-- Update users table to use string-based UUID as primary key
-- This migration converts from numeric ID + firebase_uid to string ID only

-- First, create a temporary table with the new structure
CREATE TABLE users_new (
    id VARCHAR(128) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Copy data from old table to new table, using firebase_uid as the new id
INSERT INTO users_new (id, name, email, created_at, updated_at, deleted_at)
SELECT 
    COALESCE(firebase_uid, 'user_' || id::text) as id,  -- Use firebase_uid if available, otherwise generate from old id
    name,
    email,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE firebase_uid IS NOT NULL AND firebase_uid != '';

-- Drop the old table
DROP TABLE users;

-- Rename the new table
ALTER TABLE users_new RENAME TO users;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add comment
COMMENT ON TABLE users IS 'Stores user information with string-based UUID primary key'; 