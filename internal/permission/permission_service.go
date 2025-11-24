package permission

import (
	"context"
	"errors"

	"github.com/easy-comerce/backend/internal/permission/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/option"
)

type PermissionService interface {
	CreatePermissionsIfNotExists(ctx context.Context, permissionNames []string, showDeleted *bool) error
	GetAllPermissions(ctx context.Context, sortBy string, sortOrder string, showDeleted *bool) ([]PermissionEntity, error)
	GetAllPermissionsGroup(ctx context.Context, sortBy string, sortOrder string) ([]model.PermissionGroup, error)
	GetPermissionByID(ctx context.Context, id uint) (*PermissionEntity, error)
	GetPermissionByName(ctx context.Context, name string) (*PermissionEntity, error)
	GetPermissionsByIDs(ctx context.Context, ids []uint) ([]PermissionEntity, error)
	GetUserStatusAndPermission(userID uint, permission string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
	GetUserStatusAndAnyPermission(userID uint, permissions []string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
}

type permissionService struct {
	repo PermissionRepository
}

func NewPermissionService(repo PermissionRepository) PermissionService {
	return &permissionService{
		repo: repo,
	}
}

func (s *permissionService) CreatePermissionsIfNotExists(ctx context.Context, permissionNames []string, showDeleted *bool) error {
	options := option.QueryOptions{
		ShowDeleted: showDeleted,
	}
	options.AddInFilter(c.PermissionName, permissionNames)

	existingPermissions, err := s.repo.GetAll(ctx, &options)
	if err != nil {
		return err
	}

	existingNamesMap := make(map[string]bool, len(existingPermissions))
	for _, perm := range existingPermissions {
		existingNamesMap[perm.Name] = true
	}

	var toCreatePermissions []PermissionEntity
	for _, name := range permissionNames {
		if !existingNamesMap[name] {
			value := c.PermissionMap[name]

			toCreatePermissions = append(
				toCreatePermissions,
				PermissionEntity{
					Name:        name,
					Description: &value.Desc,
					GroupName:   &value.GroupName,
				},
			)
		}
	}

	if len(toCreatePermissions) > 0 {
		return s.repo.CreateInBatch(ctx, &toCreatePermissions)
	}

	return nil
}

func (s *permissionService) GetAllPermissions(ctx context.Context, sortBy string, sortOrder string, showDeleted *bool) ([]PermissionEntity, error) {
	options := option.QueryOptions{
		ShowDeleted:    showDeleted,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		SortableFields: []string{c.PermissionName, c.PermissionDescription, c.PermissionGroupName},
	}

	return s.repo.GetAll(ctx, &options)
}

func (s *permissionService) GetAllPermissionsGroup(ctx context.Context, sortBy string, sortOrder string) ([]model.PermissionGroup, error) {
	fields := []string{
		c.PermissionName,
		c.PermissionDescription,
		c.PermissionGroupName,
	}

	options := option.QueryOptions{
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		SortableFields: fields,
	}

	permissions, err := s.repo.GetAll(ctx, &options, fields...)
	if err != nil {
		return nil, err
	}

	groupOrder := make([]string, 0, len(permissions))
	groupMap := make(map[string][]model.PermissionResponse, len(permissions))

	for _, p := range permissions {
		if p.GroupName == nil || *p.GroupName == "" {
			continue
		}

		key := *p.GroupName
		if _, exists := groupMap[key]; !exists {
			groupOrder = append(groupOrder, key)
		}

		groupMap[key] = append(groupMap[key], p.ToResponse())
	}

	result := make([]model.PermissionGroup, 0, len(groupMap))
	for _, groupName := range groupOrder {
		result = append(result, model.PermissionGroup{
			GroupName:   groupName,
			Permissions: groupMap[groupName],
		})
	}

	return result, nil
}

func (s *permissionService) GetPermissionByID(ctx context.Context, id uint) (*PermissionEntity, error) {
	options := option.QueryOptions{}
	options.AddFilter(c.FieldID+" = ?", id)
	return s.repo.GetSingleByID(ctx, id, &options)
}

func (s *permissionService) GetPermissionByName(ctx context.Context, name string) (*PermissionEntity, error) {
	options := option.QueryOptions{}
	options.AddFilter(c.PermissionName+" = ?", name)
	return s.repo.GetSingleBy(ctx, &options)
}

func (s *permissionService) GetPermissionsByIDs(ctx context.Context, ids []uint) ([]PermissionEntity, error) {
	options := option.QueryOptions{}
	options.AddInFilter(c.FieldID, ids)
	return s.repo.GetAll(ctx, &options)
}

func (s *permissionService) GetUserStatusAndPermission(userID uint, permission string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	if userID <= 0 {
		return false, false, nil, false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, false, nil, false, errors.New("permission cannot be empty")
	}

	return s.repo.GetUserStatusAndPermission(userID, permission, sqlComment...)
}

func (s *permissionService) GetUserStatusAndAnyPermission(userID uint, permissions []string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, nil, false, errors.New("invalid user ID")
	}

	if len(permissions) == 0 {
		return false, false, nil, false, errors.New("permissions cannot be empty")
	}

	return s.repo.GetUserStatusAndAnyPermission(userID, permissions, sqlComment...)
}
