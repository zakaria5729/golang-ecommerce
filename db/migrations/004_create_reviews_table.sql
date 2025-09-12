-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    product_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance (commented out - uncomment when needed)
-- CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id);
-- CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
-- CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
-- CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at);
-- CREATE INDEX IF NOT EXISTS idx_reviews_updated_at ON reviews(updated_at);
-- CREATE INDEX IF NOT EXISTS idx_reviews_product_rating ON reviews(product_id, rating);
-- CREATE INDEX IF NOT EXISTS idx_reviews_user_product ON reviews(user_id, product_id);
-- CREATE INDEX IF NOT EXISTS idx_reviews_product_created ON reviews(product_id, created_at DESC);

-- Insert some sample data
INSERT INTO reviews (product_id, user_id, rating, comment) VALUES
(1, 1, 5, 'Excellent product! Highly recommended.'),
(1, 2, 4, 'Good quality, fast delivery.'),
(1, 3, 3, 'Average product, could be better.'),
(2, 1, 5, 'Amazing! Will buy again.'),
(2, 2, 2, 'Not as expected, poor quality.'),
(2, 4, 4, 'Good value for money.'),
(3, 1, 5, 'Perfect! Exactly what I needed.'),
(3, 3, 4, 'Very satisfied with the purchase.'),
(3, 4, 5, 'Outstanding quality and service.'),
(4, 2, 3, 'Decent product, nothing special.'),
(4, 3, 4, 'Good overall, minor issues.'),
(5, 1, 5, 'Fantastic! Exceeded expectations.'),
(5, 2, 4, 'Great product, would recommend.'),
(5, 4, 3, 'Good but could be improved.')
ON CONFLICT (id) DO NOTHING;
