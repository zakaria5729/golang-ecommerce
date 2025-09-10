-- Drop materialized view and related functions
-- This migration removes the materialized view approach and returns to 5-table joins

-- Drop materialized view first (this will also drop dependent objects)
DROP MATERIALIZED VIEW IF EXISTS user_permissions CASCADE;

-- Drop triggers (in case they still exist)
DROP TRIGGER IF EXISTS trigger_user_roles_refresh ON user_roles;
DROP TRIGGER IF EXISTS trigger_role_permissions_refresh ON role_permissions;
DROP TRIGGER IF EXISTS trigger_permissions_refresh ON permissions;
DROP TRIGGER IF EXISTS trigger_roles_refresh ON roles;

-- Note: Some functions may still exist but are no longer used
-- The materialized view has been successfully dropped
