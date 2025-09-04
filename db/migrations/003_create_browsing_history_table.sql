-- Create browsing_history table
CREATE TABLE IF NOT EXISTS browsing_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    viewed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- -- Create index on user_id for better performance
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_user_id ON browsing_history(user_id);

-- -- Create index on product_id for better performance
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_product_id ON browsing_history(product_id);

-- -- Create index on viewed_at for sorting
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_viewed_at ON browsing_history(viewed_at);

-- -- Create composite index for user_id and viewed_at for efficient user history queries
-- CREATE INDEX IF NOT EXISTS idx_browsing_history_user_viewed ON browsing_history(user_id, viewed_at DESC);

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
