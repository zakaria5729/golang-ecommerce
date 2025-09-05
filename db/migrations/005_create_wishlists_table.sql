-- Create wishlists table
CREATE TABLE IF NOT EXISTS wishlists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

-- Indexes for better performance (commented out - uncomment when needed)
-- CREATE INDEX IF NOT EXISTS idx_wishlists_user_id ON wishlists(user_id);
-- CREATE INDEX IF NOT EXISTS idx_wishlists_product_id ON wishlists(product_id);
-- CREATE INDEX IF NOT EXISTS idx_wishlists_created_at ON wishlists(created_at);
-- CREATE INDEX IF NOT EXISTS idx_wishlists_updated_at ON wishlists(updated_at);
-- CREATE INDEX IF NOT EXISTS idx_wishlists_user_product ON wishlists(user_id, product_id);
-- CREATE INDEX IF NOT EXISTS idx_wishlists_user_created ON wishlists(user_id, created_at DESC);

-- Insert some sample data
INSERT INTO wishlists (user_id, product_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 1),
(2, 4),
(2, 5),
(3, 2),
(3, 3),
(3, 6),
(4, 1),
(4, 2),
(4, 4),
(5, 3),
(5, 5),
(5, 6)
ON CONFLICT (user_id, product_id) DO NOTHING;
