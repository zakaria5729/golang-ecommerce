CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role_type VARCHAR(50) NOT NULL CHECK (role_type IN ('SUPER_ADMIN', 'ADMIN', 'MANAGER', 'SELLER', 'USER')),
    role_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

INSERT INTO roles (role_name, role_type, description) VALUES
('Super Admin', 'SUPER_ADMIN', 'Platform owner with all privileges'),
('Admin', 'ADMIN', 'Administrator with management privileges'),
('Manager', 'MANAGER', 'Manager with content management privileges'),
('Seller', 'SELLER', 'Seller with product management privileges'),
('User', 'USER', 'Regular user with basic privileges')
ON CONFLICT (role_name) DO NOTHING;
