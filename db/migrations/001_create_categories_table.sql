-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN DEFAULT true,
    parent_id INTEGER REFERENCES categories(id),
    priority INTEGER DEFAULT 0,
    title VARCHAR(120) NOT NULL,
    sub_title VARCHAR(255),
    path_key VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for better performance (commented out - uncomment when needed)
-- CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);
-- CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);
-- CREATE INDEX IF NOT EXISTS idx_categories_title ON categories(title);
-- CREATE INDEX IF NOT EXISTS idx_categories_priority ON categories(priority);
-- CREATE INDEX IF NOT EXISTS idx_categories_created_at ON categories(created_at);
-- CREATE INDEX IF NOT EXISTS idx_categories_updated_at ON categories(updated_at);

-- Insert some sample data
INSERT INTO categories (title, sub_title, image_url, is_active, priority) VALUES
('Electronics', 'Electronic devices and gadgets', 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Altja_j%C3%B5gi_Lahemaal.jpg/960px-Altja_j%C3%B5gi_Lahemaal.jpg', true, 0),
('Clothing', 'Fashion and apparel', 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Altja_j%C3%B5gi_Lahemaal.jpg/960px-Altja_j%C3%B5gi_Lahemaal.jpg', true, 0),
('Books', 'Books and publications', 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Altja_j%C3%B5gi_Lahemaal.jpg/960px-Altja_j%C3%B5gi_Lahemaal.jpg', true, 0)
ON CONFLICT (id) DO NOTHING;

