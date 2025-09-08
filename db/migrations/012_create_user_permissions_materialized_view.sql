-- Create materialized view for user permissions
-- This replaces the denormalized table approach with a more efficient materialized view

-- Drop the existing user_permissions materialized view if it exists
DROP MATERIALIZED VIEW IF EXISTS user_permissions CASCADE;

-- Create materialized view for user permissions
CREATE MATERIALIZED VIEW user_permissions AS
SELECT DISTINCT 
    u.id as user_id,
    p.name as permission_name,
    r.role_name,
    r.role_type,
    CURRENT_TIMESTAMP as created_at,
    CURRENT_TIMESTAMP as updated_at
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE u.deleted_at IS NULL 
  AND r.deleted_at IS NULL 
  AND p.deleted_at IS NULL;

-- Create unique index on the materialized view
CREATE UNIQUE INDEX idx_user_permissions_unique ON user_permissions(user_id, permission_name);

-- Create indexes for fast lookups
CREATE INDEX idx_user_permissions_user_id ON user_permissions(user_id);
CREATE INDEX idx_user_permissions_permission ON user_permissions(permission_name);
CREATE INDEX idx_user_permissions_role_type ON user_permissions(role_type);
CREATE INDEX idx_user_permissions_role_name ON user_permissions(role_name);

-- Create function to refresh the materialized view
CREATE OR REPLACE FUNCTION refresh_user_permissions()
RETURNS void AS $func$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_permissions;
END;
$func$ LANGUAGE plpgsql;


-- Create triggers to automatically refresh the materialized view
-- Trigger 1: Refresh when user_roles change
CREATE OR REPLACE FUNCTION trigger_refresh_user_permissions()
RETURNS TRIGGER AS $$
BEGIN
    -- Use a background job or queue in production
    -- For now, we'll refresh immediately (consider async in production)
    PERFORM refresh_user_permissions();
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_user_roles_refresh
    AFTER INSERT OR UPDATE OR DELETE ON user_roles
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_refresh_user_permissions();

-- Trigger 2: Refresh when role_permissions change
CREATE TRIGGER trigger_role_permissions_refresh
    AFTER INSERT OR UPDATE OR DELETE ON role_permissions
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_refresh_user_permissions();

-- Trigger 3: Refresh when permissions change
CREATE TRIGGER trigger_permissions_refresh
    AFTER INSERT OR UPDATE OR DELETE ON permissions
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_refresh_user_permissions();

-- Trigger 4: Refresh when roles change
CREATE TRIGGER trigger_roles_refresh
    AFTER INSERT OR UPDATE OR DELETE ON roles
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_refresh_user_permissions();

-- Initial refresh of the materialized view
REFRESH MATERIALIZED VIEW user_permissions;

-- Grant necessary permissions
GRANT SELECT ON user_permissions TO PUBLIC;
GRANT EXECUTE ON FUNCTION refresh_user_permissions() TO PUBLIC;
