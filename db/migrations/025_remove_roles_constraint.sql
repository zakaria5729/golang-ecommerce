-- Remove the roles_role_type_check constraint from the roles table
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_role_type_check;
