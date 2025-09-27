CREATE TABLE IF NOT EXISTS brands (
  id SERIAL PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  description VARCHAR(2000),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_by INT DEFAULT NULL,
  deleted_by INT DEFAULT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- CREATE INDEX idx_brands_name ON brands(name);
-- CREATE INDEX idx_brands_created_at ON brands(created_at);
-- CREATE INDEX idx_brands_deleted_at ON brands(deleted_at);
