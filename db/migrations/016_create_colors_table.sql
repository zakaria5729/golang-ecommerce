CREATE TABLE IF NOT EXISTS colors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_by INT DEFAULT NULL,
  deleted_by INT DEFAULT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- CREATE INDEX idx_colors_name ON colors(name);
-- CREATE INDEX idx_colors_created_at ON colors(created_at);
-- CREATE INDEX idx_colors_deleted_at ON colors(deleted_at);
