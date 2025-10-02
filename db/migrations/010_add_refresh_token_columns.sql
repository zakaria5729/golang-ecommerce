-- Add refresh token columns to users table
-- This migration adds refresh token functionality for JWT token renewal

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS refresh_token VARCHAR(255),
ADD COLUMN IF NOT EXISTS refresh_token_expires TIMESTAMP;

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS purchase_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS total_spent INTEGER DEFAULT 0;


-- Create index for refresh token lookup
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);

-- Create indexes for better performance (commented out for manual application)
-- CREATE INDEX IF NOT EXISTS idx_users_refresh_token_expires ON users(refresh_token_expires);
