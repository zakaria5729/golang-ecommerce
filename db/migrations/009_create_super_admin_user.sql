-- Create Super Admin User
-- This script creates a default super admin user with hashed password
-- Password: super123@admin (hashed using bcrypt with cost 12)

INSERT INTO users (email, password, name, verified, banned, created_at, updated_at) 
VALUES (
    'super123@admin.com',
    '$2a$12$oCBNfVQst8vBHvEwQ7096uIvjfQmz2JtbQlnqzY9t4l50UVN66AXa', -- super123@admin
    'Super Admin',
    true,
    false,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING;

-- Assign Super Admin role to the created user
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'super123@admin.com' 
AND r.role_type = 'SUPER_ADMIN'
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Create indexes for better performance (commented out for manual application)
-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
-- CREATE INDEX IF NOT EXISTS idx_users_verified ON users(verified);
-- CREATE INDEX IF NOT EXISTS idx_users_banned ON users(banned);
-- CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
-- CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
-- CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
-- CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
