-- Create addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    street TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    zip_code VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    address_type VARCHAR(20) DEFAULT 'shipping' CHECK (address_type IN ('shipping', 'billing')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- -- Create index on user_id for better performance
-- CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);

-- -- Create index on address_type for filtering
-- CREATE INDEX IF NOT EXISTS idx_addresses_address_type ON addresses(address_type);

-- -- Create index on is_default for filtering
-- CREATE INDEX IF NOT EXISTS idx_addresses_is_default ON addresses(is_default);

-- Insert some sample data
INSERT INTO addresses (user_id, street, city, state, zip_code, country, is_default, address_type) VALUES
(1, '123 Main Street', 'New York', 'NY', '10001', 'United States', true, 'shipping'),
(1, '456 Business Ave', 'New York', 'NY', '10002', 'United States', false, 'billing'),
(2, '789 Oak Drive', 'Los Angeles', 'CA', '90210', 'United States', true, 'shipping'),
(2, '321 Pine Street', 'Los Angeles', 'CA', '90211', 'United States', true, 'billing')
ON CONFLICT (id) DO NOTHING;
