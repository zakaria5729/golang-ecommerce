-- Add deleted_at column to all existing tables for soft delete support

-- Add deleted_at to categories table
ALTER TABLE categories ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Add deleted_at to addresses table
ALTER TABLE addresses ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_addresses_deleted_at ON addresses(deleted_at);

-- Add deleted_at to browsing_history table
ALTER TABLE browsing_history ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_browsing_history_deleted_at ON browsing_history(deleted_at);

-- Add deleted_at to reviews table
ALTER TABLE reviews ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_reviews_deleted_at ON reviews(deleted_at);

-- Add deleted_at to wishlists table
ALTER TABLE wishlists ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_wishlists_deleted_at ON wishlists(deleted_at);
