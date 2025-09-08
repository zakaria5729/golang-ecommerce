-- Clean up unused permission functions from the database
-- These functions were removed from the migration but may still exist in the database

-- Drop unused functions if they exist
DROP FUNCTION IF EXISTS get_user_permissions(INTEGER);
DROP FUNCTION IF EXISTS user_has_permission(INTEGER, VARCHAR(255));
DROP FUNCTION IF EXISTS get_user_permissions_by_role(INTEGER, VARCHAR(50));
