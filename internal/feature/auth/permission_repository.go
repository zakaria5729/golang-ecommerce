package auth

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		db: db.GetDB(),
	}
}

func (r *PermissionRepository) GetAllPermissions(include []string, sortBy, sortOrder string) ([]Permission, error) {
	var permissions []Permission

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&permissions).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch permissions", "method", "GetAllPermissions", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return permissions, err
}

func (r *PermissionRepository) GetPermissionByID(id uint, include []string) (*Permission, error) {
	var permission Permission

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Where(constants.FieldID+" = ?", id).First(&permission).Error; err != nil {
		logger.Logger.Error("Failed to fetch permission by ID", "method", "GetPermissionByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) GetPermissionByName(name string, include []string) (*Permission, error) {
	var permission Permission

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Where(PermissionName+" = ?", name).First(&permission).Error; err != nil {
		logger.Logger.Error("Failed to fetch permission by name", "method", "GetPermissionByName", "error", err, "name", name, "include", include)
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) GetPermissionsByNames(names []string, include []string) ([]Permission, error) {
	var permissions []Permission

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Where(PermissionName+" IN ?", names).Find(&permissions).Error; err != nil {
		logger.Logger.Error("Failed to fetch permissions by names", "method", "GetPermissionsByNames", "error", err, "names", names, "include", include)
		return nil, err
	}

	return permissions, nil
}

func (r *PermissionRepository) CreatePermission(permission *Permission) error {
	err := r.db.Create(permission).Error
	if err != nil {
		logger.Logger.Error("Failed to create permission", "method", "CreatePermission", "error", err, "permission", permission)
	}
	return err
}

func (r *PermissionRepository) UpdatePermission(permission *Permission) error {
	err := r.db.Save(permission).Error
	if err != nil {
		logger.Logger.Error("Failed to update permission", "method", "UpdatePermission", "error", err, "permission", permission)
	}
	return err
}

func (r *PermissionRepository) DeletePermission(id uint) error {
	err := r.db.Delete(&Permission{}, id).Error
	if err != nil {
		logger.Logger.Error("Failed to delete permission", "method", "DeletePermission", "error", err, "id", id)
	}
	return err
}

func (r *PermissionRepository) PermissionExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&Permission{}).Where(constants.FieldID+" = ?", id).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if permission exists", "method", "PermissionExists", "error", err, "id", id)
	}
	return count > 0, err
}

func (r *PermissionRepository) PermissionExistsByName(name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&Permission{}).Where(PermissionName+" = ?", name)

	if excludeID != nil {
		query = query.Where(constants.FieldID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if permission exists by name", "method", "PermissionExistsByName", "error", err, "name", name, "excludeID", excludeID)
	}
	return count > 0, err
}

func (r *PermissionRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, PermissionName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{"description"}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
