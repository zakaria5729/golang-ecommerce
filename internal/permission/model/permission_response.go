package model

type PermissionResponse struct {
	ID          uint    `json:"id"`
	Description *string `json:"description,omitempty"`
	GroupName   *string `json:"group_name,omitempty"`
	Name        string  `json:"name"`
}

type PermissionGroup struct {
	GroupName   string               `json:"group_name"`
	Permissions []PermissionResponse `json:"permissions"`
}
