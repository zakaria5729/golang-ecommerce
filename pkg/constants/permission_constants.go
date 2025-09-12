package constants

const (
	RoleTypeSuperAdmin = "SUPER_ADMIN"
	RoleTypeAdmin      = "ADMIN"
	RoleTypeManager    = "MANAGER"
	RoleTypeSeller     = "SELLER"
	RoleTypeUser       = "USER"
)

const (
	RoleNameSuperAdmin = "Super Admin"
	RoleNameAdmin      = "Admin"
	RoleNameManager    = "Manager"
	RoleNameSeller     = "Seller"
	RoleNameUser       = "User"
)

const (
	// User permissions
	PermissionUserCreate = "user.create"
	PermissionUserRead   = "user.read"
	PermissionUserUpdate = "user.update"
	PermissionUserDelete = "user.delete"

	// Role permissions
	PermissionRoleCreate = "role.create"
	PermissionRoleRead   = "role.read"
	PermissionRoleUpdate = "role.update"
	PermissionRoleDelete = "role.delete"
	PermissionRoleAssign = "role.assign"

	// Permission permissions
	PermissionPermissionRead   = "permission.read"
	PermissionPermissionCreate = "permission.create"
	PermissionPermissionUpdate = "permission.update"
	PermissionPermissionDelete = "permission.delete"

	// Category permissions
	PermissionCategoryCreate = "category.create"
	PermissionCategoryRead   = "category.read"
	PermissionCategoryUpdate = "category.update"
	PermissionCategoryDelete = "category.delete"

	// Address permissions
	PermissionAddressCreate = "address.create"
	PermissionAddressRead   = "address.read"
	PermissionAddressUpdate = "address.update"
	PermissionAddressDelete = "address.delete"

	// Review permissions
	PermissionReviewCreate = "review.create"
	PermissionReviewRead   = "review.read"
	PermissionReviewUpdate = "review.update"
	PermissionReviewDelete = "review.delete"

	// Wishlist permissions
	PermissionWishlistCreate = "wishlist.create"
	PermissionWishlistRead   = "wishlist.read"
	PermissionWishlistUpdate = "wishlist.update"
	PermissionWishlistDelete = "wishlist.delete"

	// Browsing History permissions
	PermissionBrowsingHistoryCreate = "browsing_history.create"
	PermissionBrowsingHistoryRead   = "browsing_history.read"
	PermissionBrowsingHistoryUpdate = "browsing_history.update"
	PermissionBrowsingHistoryDelete = "browsing_history.delete"

	// System permissions
	PermissionSystemAdmin      = "system.admin"
	PermissionSystemSuperAdmin = "system.super_admin"
)
