-- Add path_key columns to support file storage paths
-- This migration adds path_key functionality for storing file paths

-- Add path_key to users table for profile images
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS path_key VARCHAR(255);

-- Add path_key to categories table for category images  
ALTER TABLE categories
ADD COLUMN IF NOT EXISTS path_key VARCHAR(255);

-- Create indexes for better performance (commented out for manual application)
-- CREATE INDEX IF NOT EXISTS idx_users_path_key ON users(path_key);
-- CREATE INDEX IF NOT EXISTS idx_categories_path_key ON categories(path_key);
