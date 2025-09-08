# Semi-Stateless Permission System Usage Guide

## 🚀 **Overview**

The Semi-Stateless Permission System provides fast permission checking using a denormalized `user_permissions` table. This eliminates the need for complex 5-table joins and provides near-instant permission validation.

## 📊 **Performance Comparison**

| Method | Database Queries | Execution Time | Complexity |
|--------|------------------|----------------|------------|
| **Before (5-table join)** | 5 joins | 10-50ms | High |
| **After (denormalized)** | 1 query | 1-5ms | Low |

## 🔧 **Usage Examples**

### **1. Using Permission Middleware**

```go
// Require specific permission
mux.HandleFunc("GET /api/users", 
    middleware.ChainAuthMiddleware(
        authMiddleware.RequireAuth,
        permissionMiddleware.RequirePermission("user.read"),
    )(handler.GetUsers))

// Require any of multiple permissions
mux.HandleFunc("POST /api/products", 
    middleware.ChainAuthMiddleware(
        authMiddleware.RequireAuth,
        permissionMiddleware.RequireAnyPermission([]string{"product.create", "product.admin"}),
    )(handler.CreateProduct))

// Require admin role
mux.HandleFunc("DELETE /api/users/{id}", 
    middleware.ChainAuthMiddleware(
        authMiddleware.RequireAuth,
        permissionMiddleware.RequireAdmin,
    )(handler.DeleteUser))
```

### **2. Checking Permissions in Handlers**

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r)
    
    // Check permission directly
    hasPermission, err := middleware.CheckPermission(userID, "user.read")
    if err != nil || !hasPermission {
        response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
        return
    }
    
    // Continue with handler logic...
}
```

### **3. API Endpoints for Permission Management**

```bash
# Get user permissions
GET /v1/permissions/users/{id}

# Check specific permission
GET /v1/permissions/users/{id}/check?permission=user.read

# Sync user permissions (admin only)
POST /v1/permissions/users/{id}/sync

# Sync all users permissions (super admin only)
POST /v1/permissions/sync-all

# Get users with specific permission
GET /v1/permissions/users?permission=user.read

# Get permission statistics (admin only)
GET /v1/permissions/stats

# Cleanup orphaned permissions (super admin only)
POST /v1/permissions/cleanup
```

## 🔄 **Automatic Synchronization**

The system automatically handles permission/role changes through database triggers:

### **When Role Permissions Change:**
```sql
-- Adding permission to role
INSERT INTO role_permissions (role_id, permission_id) VALUES (1, 5);
-- ✅ Automatically updates user_permissions table

-- Removing permission from role  
DELETE FROM role_permissions WHERE role_id = 1 AND permission_id = 5;
-- ✅ Automatically updates user_permissions table
```

### **When User Roles Change:**
```sql
-- Assigning role to user
INSERT INTO user_roles (user_id, role_id) VALUES (123, 1);
-- ✅ Automatically updates user_permissions table

-- Removing role from user
DELETE FROM user_roles WHERE user_id = 123 AND role_id = 1;
-- ✅ Automatically updates user_permissions table
```

### **When Permission Names Change:**
```sql
-- Updating permission name
UPDATE permissions SET name = 'user.read.all' WHERE name = 'user.read';
-- ✅ Automatically updates user_permissions table
```

### **When Role Names/Types Change:**
```sql
-- Updating role name/type
UPDATE roles SET role_name = 'Senior Admin', role_type = 'SENIOR_ADMIN' 
WHERE role_name = 'Admin';
-- ✅ Automatically updates user_permissions table
```

## 🛠️ **Manual Synchronization**

### **Sync Single User:**
```bash
curl -X POST "http://localhost:8080/v1/permissions/users/123/sync" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### **Sync All Users:**
```bash
curl -X POST "http://localhost:8080/v1/permissions/sync-all" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### **Cleanup Orphaned Permissions:**
```bash
curl -X POST "http://localhost:8080/v1/permissions/cleanup" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📈 **Monitoring and Statistics**

### **Get Permission Statistics:**
```bash
curl "http://localhost:8080/v1/permissions/stats" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "success": true,
  "data": {
    "permission_stats": {
      "user.read": 150,
      "user.create": 25,
      "product.read": 200,
      "product.create": 10
    },
    "total_permissions": 4
  }
}
```

## 🔍 **Troubleshooting**

### **Permission Not Working:**
1. Check if user has the role: `GET /v1/permissions/users/{id}`
2. Sync user permissions: `POST /v1/permissions/users/{id}/sync`
3. Check permission statistics: `GET /v1/permissions/stats`

### **Performance Issues:**
1. Run cleanup: `POST /v1/permissions/cleanup`
2. Check database indexes are created
3. Monitor query execution times

### **Data Inconsistency:**
1. Sync all users: `POST /v1/permissions/sync-all`
2. Check for orphaned permissions: `POST /v1/permissions/cleanup`
3. Verify database triggers are working

## 🚀 **Migration Steps**

1. **Run Migration:**
   ```bash
   go run cmd/migrate/main.go
   ```

2. **Verify Data:**
   ```sql
   SELECT COUNT(*) FROM user_permissions;
   SELECT COUNT(*) FROM users u 
   JOIN user_roles ur ON u.id = ur.user_id
   JOIN roles r ON ur.role_id = r.id
   JOIN role_permissions rp ON r.id = rp.role_id
   JOIN permissions p ON rp.permission_id = p.id;
   ```

3. **Test Permission Checks:**
   ```bash
   # Test permission check
   curl "http://localhost:8080/v1/permissions/users/1/check?permission=user.read" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## 🎯 **Best Practices**

1. **Use Middleware:** Always use permission middleware for route protection
2. **Monitor Performance:** Regularly check permission statistics
3. **Cleanup Regularly:** Run cleanup operations periodically
4. **Test Changes:** Always test permission changes in development
5. **Backup Data:** Backup before major permission changes

## 🔐 **Security Considerations**

1. **Admin Endpoints:** Only super admins can sync all users
2. **Audit Logs:** All permission changes are logged
3. **Validation:** All inputs are validated
4. **Rate Limiting:** Consider rate limiting for sync operations
5. **Monitoring:** Monitor for unusual permission patterns
