package constants

const (
	RoleTypeSuperAdmin string = "SUPER_ADMIN"
	RoleTypeAdmin      string = "ADMIN"
	RoleTypeMaintainer string = "MAINTAINER"
	RoleTypeSeller     string = "SELLER"
	RoleTypeUser       string = "USER"
)

const (
	RoleNameSuperAdmin string = "Super Admin"
	RoleNameAdmin      string = "Admin"
	RoleNameMaintainer string = "Maintainer"
	RoleNameSeller     string = "Seller"
	RoleNameUser       string = "User"
)

const (
	PermissionGeneralUser      string = "general.user"
	PermissionPermissionRead   string = "permission.read"
	PermissionProductStatsRead string = "product_stats.read"

	PermissionUserCreate     string = "user.create"
	PermissionUserRead       string = "user.read"
	PermissionUserUpdate     string = "user.update"
	PermissionUserDelete     string = "user.delete"
	PermissionUserUndoDelete string = "user.undo_delete"

	PermissionRoleCreate           string = "role.create"
	PermissionRoleRead             string = "role.read"
	PermissionRoleUpdate           string = "role.update"
	PermissionRoleDelete           string = "role.delete"
	PermissionRoleAssign           string = "role.assign"
	PermissionRoleToAddPermissions string = "role.to_add_permissions"
	PermissionRoleUndoDelete       string = "role.undo_delete"

	PermissionCategoryCreate     string = "category.create"
	PermissionCategoryRead       string = "category.read"
	PermissionCategoryUpdate     string = "category.update"
	PermissionCategoryDelete     string = "category.delete"
	PermissionCategoryUndoDelete string = "category.undo_delete"

	PermissionAddressRead       string = "address.read"
	PermissionAddressDelete     string = "address.delete"
	PermissionAddressUndoDelete string = "address.undo_delete"

	PermissionReviewRead       string = "review.read"
	PermissionReviewUpdate     string = "review.update"
	PermissionReviewDelete     string = "review.delete"
	PermissionReviewUndoDelete string = "review.undo_delete"

	PermissionWishlistRead       string = "wishlist.read"
	PermissionWishlistDelete     string = "wishlist.delete"
	PermissionWishlistUndoDelete string = "wishlist.undo_delete"

	PermissionBrandCreate     string = "brand.create"
	PermissionBrandRead       string = "brand.read"
	PermissionBrandUpdate     string = "brand.update"
	PermissionBrandDelete     string = "brand.delete"
	PermissionBrandUndoDelete string = "brand.undo_delete"

	PermissionColorCreate     string = "color.create"
	PermissionColorRead       string = "color.read"
	PermissionColorUpdate     string = "color.update"
	PermissionColorDelete     string = "color.delete"
	PermissionColorUndoDelete string = "color.undo_delete"

	PermissionAttributeTypeCreate     string = "attribute_type.create"
	PermissionAttributeTypeRead       string = "attribute_type.read"
	PermissionAttributeTypeUpdate     string = "attribute_type.update"
	PermissionAttributeTypeDelete     string = "attribute_type.delete"
	PermissionAttributeTypeUndoDelete string = "attribute_type.undo_delete"

	PermissionAttributeOptionCreate     string = "attribute_option.create"
	PermissionAttributeOptionRead       string = "attribute_option.read"
	PermissionAttributeOptionUpdate     string = "attribute_option.update"
	PermissionAttributeOptionDelete     string = "attribute_option.delete"
	PermissionAttributeOptionUndoDelete string = "attribute_option.undo_delete"

	PermissionSizeCategoryCreate     string = "size_category.create"
	PermissionSizeCategoryRead       string = "size_category.read"
	PermissionSizeCategoryUpdate     string = "size_category.update"
	PermissionSizeCategoryDelete     string = "size_category.delete"
	PermissionSizeCategoryUndoDelete string = "size_category.undo_delete"

	PermissionSizeOptionCreate     string = "size_option.create"
	PermissionSizeOptionRead       string = "size_option.read"
	PermissionSizeOptionUpdate     string = "size_option.update"
	PermissionSizeOptionDelete     string = "size_option.delete"
	PermissionSizeOptionUndoDelete string = "size_option.undo_delete"
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
