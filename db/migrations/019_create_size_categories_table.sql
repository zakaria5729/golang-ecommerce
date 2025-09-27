CREATE TABLE IF NOT EXISTS size_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_by INT DEFAULT NULL,
  deleted_by INT DEFAULT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- CREATE INDEX idx_size_categories_name ON size_categories (name);
-- CREATE INDEX idx_size_categories_created_at ON size_categories (created_at);
-- CREATE INDEX idx_size_categories_updated_at ON size_categories (updated_at);
