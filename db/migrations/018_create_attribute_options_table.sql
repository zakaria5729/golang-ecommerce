CREATE TABLE IF NOT EXISTS attribute_options (
  id SERIAL PRIMARY KEY,
  attribute_type_id INT DEFAULT NULL,
  attribute_option_name VARCHAR(100) DEFAULT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_by INT DEFAULT NULL,
  deleted_by INT DEFAULT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT fk_attroption_attr_type_id FOREIGN KEY (attribute_type_id) REFERENCES attribute_types (id) ON DELETE RESTRICT
);

-- CREATE INDEX idx_attribute_options_attribute_type_id ON attribute_options(attribute_type_id);
-- CREATE INDEX idx_attribute_options_attribute_option_name ON attribute_options(attribute_option_name);
-- CREATE INDEX idx_attribute_options_created_at ON attribute_options(created_at);
-- CREATE INDEX idx_attribute_options_deleted_at ON attribute_options(deleted_at);
