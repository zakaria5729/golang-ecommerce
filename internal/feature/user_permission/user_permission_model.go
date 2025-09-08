package user_permission

type UserPermission struct {
	UserID         uint   `json:"user_id" gorm:"not null;column:user_id"`
	PermissionName string `json:"permission_name" gorm:"not null;column:permission_name"`
	RoleName       string `json:"role_name" gorm:"not null;column:role_name"`
	RoleType       string `json:"role_type" gorm:"not null;column:role_type"`
}

const (
	UserPermissionUserID         = "user_id"
	UserPermissionPermissionName = "permission_name"
	UserPermissionRoleName       = "role_name"
	UserPermissionRoleType       = "role_type"
)
