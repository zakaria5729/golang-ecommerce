-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(120) NOT NULL,
    sub_title VARCHAR(255),
    image_url TEXT,
    parent_id INTEGER REFERENCES categories(id),
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on parent_id for better performance
-- CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- Create index on is_active for filtering
-- CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);

-- Insert some sample data
INSERT INTO categories (title, sub_title, image_url, is_active, priority) VALUES
('Electronics', 'Electronic devices and gadgets', 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Altja_j%C3%B5gi_Lahemaal.jpg/960px-Altja_j%C3%B5gi_Lahemaal.jpg', true, 0),
('Clothing', 'Fashion and apparel', 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Altja_j%C3%B5gi_Lahemaal.jpg/960px-Altja_j%C3%B5gi_Lahemaal.jpg', true, 0),
('Books', 'Books and publications', 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Altja_j%C3%B5gi_Lahemaal.jpg/960px-Altja_j%C3%B5gi_Lahemaal.jpg', true, 0)
ON CONFLICT (id) DO NOTHING;

