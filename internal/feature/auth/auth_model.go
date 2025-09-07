package auth

// This file contains the junction table models for many-to-many relationships

type UserRole struct {
	UserID uint `gorm:"primaryKey;column:user_id"`
	RoleID uint `gorm:"primaryKey;column:role_id"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;column:role_id"`
	PermissionID uint `gorm:"primaryKey;column:permission_id"`
}
