CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    banned BOOLEAN DEFAULT FALSE,
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id 
FROM roles r, permissions p 
WHERE r.role_type = 'SUPER_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id 
FROM roles r, permissions p 
WHERE r.role_type = 'ADMIN' 
AND p.name IN (
    'user.create', 'user.read', 'user.update', 'user.delete',
    'role.read', 'permission.read',
    'category.create', 'category.read', 'category.update', 'category.delete',
    'address.create', 'address.read', 'address.update', 'address.delete',
    'review.create', 'review.read', 'review.update', 'review.delete',
    'wishlist.create', 'wishlist.read', 'wishlist.update', 'wishlist.delete',
    'browsing_history.create', 'browsing_history.read', 'browsing_history.update', 'browsing_history.delete',
    'system.admin'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id 
FROM roles r, permissions p 
WHERE r.role_type = 'MANAGER' 
AND p.name IN (
    'user.read', 'user.update',
    'category.create', 'category.read', 'category.update', 'category.delete',
    'address.create', 'address.read', 'address.update', 'address.delete',
    'review.create', 'review.read', 'review.update', 'review.delete',
    'wishlist.create', 'wishlist.read', 'wishlist.update', 'wishlist.delete',
    'browsing_history.create', 'browsing_history.read', 'browsing_history.update', 'browsing_history.delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id 
FROM roles r, permissions p 
WHERE r.role_type = 'SELLER' 
AND p.name IN (
    'category.read',
    'address.create', 'address.read', 'address.update', 'address.delete',
    'review.create', 'review.read', 'review.update',
    'wishlist.create', 'wishlist.read', 'wishlist.update', 'wishlist.delete',
    'browsing_history.create', 'browsing_history.read', 'browsing_history.update', 'browsing_history.delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id 
FROM roles r, permissions p 
WHERE r.role_type = 'USER' 
AND p.name IN (
    'address.create', 'address.read', 'address.update', 'address.delete',
    'review.create', 'review.read', 'review.update',
    'wishlist.create', 'wishlist.read', 'wishlist.update', 'wishlist.delete',
    'browsing_history.create', 'browsing_history.read', 'browsing_history.update', 'browsing_history.delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;
