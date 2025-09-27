CREATE TABLE IF NOT EXISTS attribute_types (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) DEFAULT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_by INT DEFAULT NULL,
  deleted_by INT DEFAULT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- CREATE INDEX idx_attribute_types_name ON attribute_types(name);
-- CREATE INDEX idx_attribute_types_created_at ON attribute_types(created_at);
-- CREATE INDEX idx_attribute_types_deleted_at ON attribute_types(deleted_at);
