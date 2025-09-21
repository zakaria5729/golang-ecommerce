package constants

const (
	RoleTypeSuperAdmin = "SUPER_ADMIN"
	RoleTypeAdmin      = "ADMIN"
	RoleTypeMaintainer = "MAINTAINER"
	RoleTypeSeller     = "SELLER"
	RoleTypeUser       = "USER"
)

const (
	RoleNameSuperAdmin = "Super Admin"
	RoleNameAdmin      = "Admin"
	RoleNameMaintainer = "Maintainer"
	RoleNameSeller     = "Seller"
	RoleNameUser       = "User"
)

const (
	PermissionUserCreate     = "user.create"
	PermissionUserRead       = "user.read"
	PermissionUserUpdate     = "user.update"
	PermissionUserDelete     = "user.delete"
	PermissionUserUndoDelete = "user.undo_delete"

	PermissionRoleCreate           = "role.create"
	PermissionRoleRead             = "role.read"
	PermissionRoleUpdate           = "role.update"
	PermissionRoleDelete           = "role.delete"
	PermissionRoleAssign           = "role.assign"
	PermissionRoleToAddPermissions = "role.to_add_permissions"
	PermissionRoleUndoDelete       = "role.undo_delete"

	PermissionPermissionRead   = "permission.read"
	PermissionProductStatsRead = "product_stats.read"

	PermissionCategoryCreate     = "category.create"
	PermissionCategoryRead       = "category.read"
	PermissionCategoryUpdate     = "category.update"
	PermissionCategoryDelete     = "category.delete"
	PermissionCategoryUndoDelete = "category.undo_delete"

	PermissionAddressRead       = "address.read"
	PermissionAddressDelete     = "address.delete"
	PermissionAddressUndoDelete = "address.undo_delete"

	PermissionReviewRead       = "review.read"
	PermissionReviewDelete     = "review.delete"
	PermissionReviewUndoDelete = "review.undo_delete"

	PermissionWishlistRead       = "wishlist.read"
	PermissionWishlistDelete     = "wishlist.delete"
	PermissionWishlistUndoDelete = "wishlist.undo_delete"
)

var PermissionMap = map[string]string{
	PermissionUserCreate:     "Create User",
	PermissionUserRead:       "Read User",
	PermissionUserUpdate:     "Update User",
	PermissionUserDelete:     "Delete User",
	PermissionUserUndoDelete: "Undo Delete User",

	PermissionRoleCreate:           "Create Role",
	PermissionRoleRead:             "Read Role",
	PermissionRoleUpdate:           "Update Role",
	PermissionRoleDelete:           "Delete Role",
	PermissionRoleAssign:           "Assign Role",
	PermissionRoleToAddPermissions: "Add Permissions to Role",
	PermissionRoleUndoDelete:       "Undo Delete Role",

	PermissionPermissionRead:   "Read Permission",
	PermissionProductStatsRead: "Read Product Stats",

	PermissionCategoryCreate:     "Create Category",
	PermissionCategoryRead:       "Read Category",
	PermissionCategoryUpdate:     "Update Category",
	PermissionCategoryDelete:     "Delete Category",
	PermissionCategoryUndoDelete: "Undo Delete Category",

	PermissionAddressRead:       "Read Address",
	PermissionAddressDelete:     "Delete Address",
	PermissionAddressUndoDelete: "Undo Delete Address",

	PermissionReviewRead:       "Read Review",
	PermissionReviewDelete:     "Delete Review",
	PermissionReviewUndoDelete: "Undo Delete Review",

	PermissionWishlistRead:       "Read Wishlist",
	PermissionWishlistDelete:     "Delete Wishlist",
	PermissionWishlistUndoDelete: "Undo Delete Wishlist",
}
