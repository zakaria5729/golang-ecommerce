-- Create addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    is_default BOOLEAN DEFAULT FALSE,
    user_id INTEGER NOT NULL,
    address_type VARCHAR(20) DEFAULT 'shipping' CHECK (address_type IN ('shipping', 'billing')),
    zip_code VARCHAR(20),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL,
    street TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance (commented out - uncomment when needed)
-- CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);
-- CREATE INDEX IF NOT EXISTS idx_addresses_address_type ON addresses(address_type);
-- CREATE INDEX IF NOT EXISTS idx_addresses_is_default ON addresses(is_default);
-- CREATE INDEX IF NOT EXISTS idx_addresses_city ON addresses(city);
-- CREATE INDEX IF NOT EXISTS idx_addresses_country ON addresses(country);
-- CREATE INDEX IF NOT EXISTS idx_addresses_created_at ON addresses(created_at);
-- CREATE INDEX IF NOT EXISTS idx_addresses_updated_at ON addresses(updated_at);
-- CREATE INDEX IF NOT EXISTS idx_addresses_user_id_address_type ON addresses(user_id, address_type);
-- CREATE INDEX IF NOT EXISTS idx_addresses_user_id_is_default ON addresses(user_id, is_default);

-- Insert some sample data
INSERT INTO addresses (user_id, street, city, state, zip_code, country, is_default, address_type) VALUES
(1, '123 Main Street', 'New York', 'NY', '10001', 'United States', true, 'shipping'),
(1, '456 Business Ave', 'New York', 'NY', '10002', 'United States', false, 'billing'),
(2, '789 Oak Drive', 'Los Angeles', 'CA', '90210', 'United States', true, 'shipping'),
(2, '321 Pine Street', 'Los Angeles', 'CA', '90211', 'United States', true, 'billing')
ON CONFLICT (id) DO NOTHING;
