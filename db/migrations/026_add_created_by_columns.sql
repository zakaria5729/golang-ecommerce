-- Add created_by column to all tables that need it

-- Users table
ALTER TABLE users ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE users ALTER COLUMN created_by SET DEFAULT NULL;

-- Roles table
ALTER TABLE roles ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE roles ALTER COLUMN created_by SET DEFAULT NULL;

-- Categories table
ALTER TABLE categories ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE categories ALTER COLUMN created_by SET DEFAULT NULL;

-- Addresses table
ALTER TABLE addresses ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE addresses ALTER COLUMN created_by SET DEFAULT NULL;

-- Brands table
ALTER TABLE brands ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE brands ALTER COLUMN created_by SET DEFAULT NULL;

-- Colors table
ALTER TABLE colors ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE colors ALTER COLUMN created_by SET DEFAULT NULL;

-- Products table
ALTER TABLE products ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE products ALTER COLUMN created_by SET DEFAULT NULL;

-- Product_stats table
ALTER TABLE product_stats ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE product_stats ALTER COLUMN created_by SET DEFAULT NULL;

-- Browsing_histories table
ALTER TABLE browsing_histories ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE browsing_histories ALTER COLUMN created_by SET DEFAULT NULL;

-- Reviews table
ALTER TABLE reviews ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE reviews ALTER COLUMN created_by SET DEFAULT NULL;

-- Wishlists table
ALTER TABLE wishlists ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE wishlists ALTER COLUMN created_by SET DEFAULT NULL;

-- Attribute_types table
ALTER TABLE attribute_types ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE attribute_types ALTER COLUMN created_by SET DEFAULT NULL;
