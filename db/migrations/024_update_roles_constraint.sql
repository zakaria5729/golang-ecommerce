-- Update the roles table constraint to include MAINTAINER role type
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_role_type_check;
ALTER TABLE roles ADD CONSTRAINT roles_role_type_check CHECK (role_type IN ('SUPER_ADMIN', 'ADMIN', 'MANAGER', 'MAINTAINER', 'SELLER', 'USER'));

-- Insert MAINTAINER role if it doesn't exist
INSERT INTO roles (role_name, role_type, description) VALUES
('Maintainer', 'MAINTAINER', 'Maintainer with system maintenance privileges')
ON CONFLICT (role_name) DO NOTHING;
