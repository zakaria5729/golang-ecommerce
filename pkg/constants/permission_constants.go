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
	PermissionGeneralUser = "general.user"

	PermissionPermissionRead   = "permission.read"
	PermissionProductStatsRead = "product_stats.read"

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

	PermissionBrandCreate     = "brand.create"
	PermissionBrandRead       = "brand.read"
	PermissionBrandUpdate     = "brand.update"
	PermissionBrandDelete     = "brand.delete"
	PermissionBrandUndoDelete = "brand.undo_delete"

	PermissionColorCreate     = "color.create"
	PermissionColorRead       = "color.read"
	PermissionColorUpdate     = "color.update"
	PermissionColorDelete     = "color.delete"
	PermissionColorUndoDelete = "color.undo_delete"

	PermissionAttributeTypeCreate     = "attribute_type.create"
	PermissionAttributeTypeRead       = "attribute_type.read"
	PermissionAttributeTypeUpdate     = "attribute_type.update"
	PermissionAttributeTypeDelete     = "attribute_type.delete"
	PermissionAttributeTypeUndoDelete = "attribute_type.undo_delete"

	PermissionAttributeOptionCreate     = "attribute_option.create"
	PermissionAttributeOptionRead       = "attribute_option.read"
	PermissionAttributeOptionUpdate     = "attribute_option.update"
	PermissionAttributeOptionDelete     = "attribute_option.delete"
	PermissionAttributeOptionUndoDelete = "attribute_option.undo_delete"

	PermissionSizeCategoryCreate     = "size_category.create"
	PermissionSizeCategoryRead       = "size_category.read"
	PermissionSizeCategoryUpdate     = "size_category.update"
	PermissionSizeCategoryDelete     = "size_category.delete"
	PermissionSizeCategoryUndoDelete = "size_category.undo_delete"

	PermissionSizeOptionCreate     = "size_option.create"
	PermissionSizeOptionRead       = "size_option.read"
	PermissionSizeOptionUpdate     = "size_option.update"
	PermissionSizeOptionDelete     = "size_option.delete"
	PermissionSizeOptionUndoDelete = "size_option.undo_delete"
)

var PermissionMap = map[string]string{
	PermissionGeneralUser: "General User",

	PermissionPermissionRead:   "Read Permission",
	PermissionProductStatsRead: "Read Product Stats",

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

	PermissionBrandCreate:     "Create Brand",
	PermissionBrandRead:       "Read Brand",
	PermissionBrandUpdate:     "Update Brand",
	PermissionBrandDelete:     "Delete Brand",
	PermissionBrandUndoDelete: "Undo Delete Brand",

	PermissionColorCreate:     "Create Color",
	PermissionColorRead:       "Read Color",
	PermissionColorUpdate:     "Update Color",
	PermissionColorDelete:     "Delete Color",
	PermissionColorUndoDelete: "Undo Delete Color",

	PermissionAttributeTypeCreate:     "Create Attribute Type",
	PermissionAttributeTypeRead:       "Read Attribute Type",
	PermissionAttributeTypeUpdate:     "Update Attribute Type",
	PermissionAttributeTypeDelete:     "Delete Attribute Type",
	PermissionAttributeTypeUndoDelete: "Undo Delete Attribute Type",

	PermissionAttributeOptionCreate:     "Create Attribute Option",
	PermissionAttributeOptionRead:       "Read Attribute Option",
	PermissionAttributeOptionUpdate:     "Update Attribute Option",
	PermissionAttributeOptionDelete:     "Delete Attribute Option",
	PermissionAttributeOptionUndoDelete: "Undo Delete Attribute Option",

	PermissionSizeCategoryCreate:     "Create Size Category",
	PermissionSizeCategoryRead:       "Read Size Category",
	PermissionSizeCategoryUpdate:     "Update Size Category",
	PermissionSizeCategoryDelete:     "Delete Size Category",
	PermissionSizeCategoryUndoDelete: "Undo Delete Size Category",

	PermissionSizeOptionCreate:     "Create Size Option",
	PermissionSizeOptionRead:       "Read Size Option",
	PermissionSizeOptionUpdate:     "Update Size Option",
	PermissionSizeOptionDelete:     "Delete Size Option",
	PermissionSizeOptionUndoDelete: "Undo Delete Size Option",
}
