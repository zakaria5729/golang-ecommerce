-- Add created_by column to all tables that don't have it

-- Users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Roles table
ALTER TABLE roles ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Categories table
ALTER TABLE categories ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Addresses table
ALTER TABLE addresses ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Brands table
ALTER TABLE brands ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Colors table
ALTER TABLE colors ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Product_stats table
ALTER TABLE product_stats ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Browsing_histories table
ALTER TABLE browsing_histories ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Reviews table
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Wishlists table
ALTER TABLE wishlists ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Attribute_types table
ALTER TABLE attribute_types ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Size_options table
ALTER TABLE size_options ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Size_categories table
ALTER TABLE size_categories ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Note: user_roles and role_permissions are junction tables and typically don't need created_by
-- as they represent many-to-many relationships rather than user-created entities

-- Size_options table
ALTER TABLE size_options ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

-- Size_categories table
ALTER TABLE size_categories ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;

ALTER TABLE permissions ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS created_by INT DEFAULT NULL;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS updated_by INT DEFAULT NULL;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS deleted_by INT DEFAULT NULL;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;