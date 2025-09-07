-- Add refresh token columns to users table
-- This migration adds refresh token functionality for JWT token renewal

ALTER TABLE users 
ADD COLUMN refresh_token VARCHAR(255),
ADD COLUMN refresh_token_expires TIMESTAMP;

-- Create index for refresh token lookup
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);

-- Create indexes for better performance (commented out for manual application)
-- CREATE INDEX IF NOT EXISTS idx_users_refresh_token_expires ON users(refresh_token_expires);
