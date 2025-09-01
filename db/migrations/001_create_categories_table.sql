-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    sub_title VARCHAR(255),
    image_url TEXT,
    parent_id INTEGER REFERENCES categories(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on parent_id for better performance
-- CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- Create index on is_active for filtering
-- CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);

-- Insert some sample data
INSERT INTO categories (title, sub_title, image_url, is_active) VALUES
('Electronics', 'Electronic devices and gadgets', 'https://example.com/electronics.jpg', true),
('Clothing', 'Fashion and apparel', 'https://example.com/clothing.jpg', true),
('Books', 'Books and publications', 'https://example.com/books.jpg', true)
ON CONFLICT (id) DO NOTHING;
