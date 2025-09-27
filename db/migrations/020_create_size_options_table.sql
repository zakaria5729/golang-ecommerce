CREATE TABLE IF NOT EXISTS size_options (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  sort_order INT DEFAULT NULL,
  size_category_id INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_by INT DEFAULT NULL,
  deleted_by INT DEFAULT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT fk_sizeoption_size_category FOREIGN KEY (size_category_id) REFERENCES size_categories (id) ON DELETE RESTRICT
);

-- CREATE INDEX idx_size_options_name ON size_options (name);
-- CREATE INDEX idx_size_options_sort_order ON size_options (sort_order);
-- CREATE INDEX idx_size_options_size_category_id ON size_options (size_category_id);
-- CREATE INDEX idx_size_options_created_at ON size_options (created_at);
-- CREATE INDEX idx_size_options_updated_at ON size_options (updated_at);
