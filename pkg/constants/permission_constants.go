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

type PermissionDetails struct {
	Desc      string
	GroupName string
}

var PermissionMap = map[string]PermissionDetails{
	// General
	PermissionGeneralUser: {
		Desc:      "General User",
		GroupName: "General",
	},

	// User
	PermissionUserCreate: {
		Desc:      "Create User",
		GroupName: "User",
	},
	PermissionUserRead: {
		Desc:      "Read User",
		GroupName: "User",
	},
	PermissionUserUpdate: {
		Desc:      "Update User",
		GroupName: "User",
	},
	PermissionUserDelete: {
		Desc:      "Delete User",
		GroupName: "User",
	},
	PermissionUserUndoDelete: {
		Desc:      "Undo Delete User",
		GroupName: "User",
	},

	// Permission
	PermissionPermissionRead: {
		Desc:      "Read Permission",
		GroupName: "Permission",
	},
	PermissionProductStatsRead: {
		Desc:      "Read Product Stats",
		GroupName: "Product Stats",
	},

	// Role
	PermissionRoleCreate: {
		Desc:      "Create Role",
		GroupName: "Role",
	},
	PermissionRoleRead: {
		Desc:      "Read Role",
		GroupName: "Role",
	},
	PermissionRoleUpdate: {
		Desc:      "Update Role",
		GroupName: "Role",
	},
	PermissionRoleDelete: {
		Desc:      "Delete Role",
		GroupName: "Role",
	},
	PermissionRoleAssign: {
		Desc:      "Assign Role",
		GroupName: "Role",
	},
	PermissionRoleToAddPermissions: {
		Desc:      "Add Permissions to Role",
		GroupName: "Role",
	},
	PermissionRoleUndoDelete: {
		Desc:      "Undo Delete Role",
		GroupName: "Role",
	},

	// Category
	PermissionCategoryCreate: {
		Desc:      "Create Category",
		GroupName: "Category",
	},
	PermissionCategoryRead: {
		Desc:      "Read Category",
		GroupName: "Category",
	},
	PermissionCategoryUpdate: {
		Desc:      "Update Category",
		GroupName: "Category",
	},
	PermissionCategoryDelete: {
		Desc:      "Delete Category",
		GroupName: "Category",
	},
	PermissionCategoryUndoDelete: {
		Desc:      "Undo Delete Category",
		GroupName: "Category",
	},

	// Address
	PermissionAddressRead: {
		Desc:      "Read Address",
		GroupName: "Address",
	},
	PermissionAddressDelete: {
		Desc:      "Delete Address",
		GroupName: "Address",
	},
	PermissionAddressUndoDelete: {
		Desc:      "Undo Delete Address",
		GroupName: "Address",
	},

	// Review
	PermissionReviewRead: {
		Desc:      "Read Review",
		GroupName: "Review",
	},
	PermissionReviewDelete: {
		Desc:      "Delete Review",
		GroupName: "Review",
	},
	PermissionReviewUndoDelete: {
		Desc:      "Undo Delete Review",
		GroupName: "Review",
	},

	// Wishlist
	PermissionWishlistRead: {
		Desc:      "Read Wishlist",
		GroupName: "Wishlist",
	},
	PermissionWishlistDelete: {
		Desc:      "Delete Wishlist",
		GroupName: "Wishlist",
	},
	PermissionWishlistUndoDelete: {
		Desc:      "Undo Delete Wishlist",
		GroupName: "Wishlist",
	},

	// Brand
	PermissionBrandCreate: {
		Desc:      "Create Brand",
		GroupName: "Brand",
	},
	PermissionBrandRead: {
		Desc:      "Read Brand",
		GroupName: "Brand",
	},
	PermissionBrandUpdate: {
		Desc:      "Update Brand",
		GroupName: "Brand",
	},
	PermissionBrandDelete: {
		Desc:      "Delete Brand",
		GroupName: "Brand",
	},
	PermissionBrandUndoDelete: {
		Desc:      "Undo Delete Brand",
		GroupName: "Brand",
	},

	// Color
	PermissionColorCreate: {
		Desc:      "Create Color",
		GroupName: "Color",
	},
	PermissionColorRead: {
		Desc:      "Read Color",
		GroupName: "Color",
	},
	PermissionColorUpdate: {
		Desc:      "Update Color",
		GroupName: "Color",
	},
	PermissionColorDelete: {
		Desc:      "Delete Color",
		GroupName: "Color",
	},
	PermissionColorUndoDelete: {
		Desc:      "Undo Delete Color",
		GroupName: "Color",
	},

	// Attribute Type
	PermissionAttributeTypeCreate: {
		Desc:      "Create Attribute Type",
		GroupName: "Attribute Type",
	},
	PermissionAttributeTypeRead: {
		Desc:      "Read Attribute Type",
		GroupName: "Attribute Type",
	},
	PermissionAttributeTypeUpdate: {
		Desc:      "Update Attribute Type",
		GroupName: "Attribute Type",
	},
	PermissionAttributeTypeDelete: {
		Desc:      "Delete Attribute Type",
		GroupName: "Attribute Type",
	},
	PermissionAttributeTypeUndoDelete: {
		Desc:      "Undo Delete Attribute Type",
		GroupName: "Attribute Type",
	},

	// Attribute Option
	PermissionAttributeOptionCreate: {
		Desc:      "Create Attribute Option",
		GroupName: "Attribute Option",
	},
	PermissionAttributeOptionRead: {
		Desc:      "Read Attribute Option",
		GroupName: "Attribute Option",
	},
	PermissionAttributeOptionUpdate: {
		Desc:      "Update Attribute Option",
		GroupName: "Attribute Option",
	},
	PermissionAttributeOptionDelete: {
		Desc:      "Delete Attribute Option",
		GroupName: "Attribute Option",
	},
	PermissionAttributeOptionUndoDelete: {
		Desc:      "Undo Delete Attribute Option",
		GroupName: "Attribute Option",
	},

	// Size Category
	PermissionSizeCategoryCreate: {
		Desc:      "Create Size Category",
		GroupName: "Size Category",
	},
	PermissionSizeCategoryRead: {
		Desc:      "Read Size Category",
		GroupName: "Size Category",
	},
	PermissionSizeCategoryUpdate: {
		Desc:      "Update Size Category",
		GroupName: "Size Category",
	},
	PermissionSizeCategoryDelete: {
		Desc:      "Delete Size Category",
		GroupName: "Size Category",
	},
	PermissionSizeCategoryUndoDelete: {
		Desc:      "Undo Delete Size Category",
		GroupName: "Size Category",
	},

	// Size Option
	PermissionSizeOptionCreate: {
		Desc:      "Create Size Option",
		GroupName: "Size Option",
	},
	PermissionSizeOptionRead: {
		Desc:      "Read Size Option",
		GroupName: "Size Option",
	},
	PermissionSizeOptionUpdate: {
		Desc:      "Update Size Option",
		GroupName: "Size Option",
	},
	PermissionSizeOptionDelete: {
		Desc:      "Delete Size Option",
		GroupName: "Size Option",
	},
	PermissionSizeOptionUndoDelete: {
		Desc:      "Undo Delete Size Option",
		GroupName: "Size Option",
	},
}
