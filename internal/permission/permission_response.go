package permission

type PermissionGroup struct {
	GroupName   string       `json:"group_name"`
	Permissions []Permission `json:"permissions"`
}
