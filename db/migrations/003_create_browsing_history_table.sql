-- Create browsing_history table
CREATE TABLE IF NOT EXISTS browsing_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    viewed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance (commented out - uncomment when needed)
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_user_id ON browsing_history(user_id);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_product_id ON browsing_history(product_id);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_viewed_at ON browsing_history(viewed_at);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_created_at ON browsing_history(created_at);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_updated_at ON browsing_history(updated_at);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_user_viewed ON browsing_history(user_id, viewed_at DESC);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_product_viewed ON browsing_history(product_id, viewed_at DESC);
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_user_product ON browsing_history(user_id, product_id);

-- Insert some sample data
INSERT INTO browsing_history (user_id, product_id, viewed_at) VALUES
(1, 1, '2025-09-04 10:00:00+06:00'),
(1, 2, '2025-09-04 10:15:00+06:00'),
(1, 3, '2025-09-04 10:30:00+06:00'),
(1, 1, '2025-09-04 11:00:00+06:00'),
(2, 2, '2025-09-04 11:15:00+06:00'),
(2, 4, '2025-09-04 11:30:00+06:00'),
(1, 5, '2025-09-04 12:00:00+06:00'),
(2, 1, '2025-09-04 12:15:00+06:00')
ON CONFLICT (id) DO NOTHING;
