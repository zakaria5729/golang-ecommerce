CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role_type VARCHAR(50) NOT NULL,
    role_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO roles (role_name, role_type, description) VALUES
('Super Admin', 'SUPER_ADMIN', 'Platform owner with all privileges'),
('Admin', 'ADMIN', 'Administrator with management privileges'),
('Manager', 'MANAGER', 'Manager with content management privileges'),
('Seller', 'SELLER', 'Seller with product management privileges'),
('User', 'USER', 'Regular user with basic privileges')
ON CONFLICT (role_name) DO NOTHING;
