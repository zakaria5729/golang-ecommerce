// package shared

// // UserRepositoryInterface defines the interface for user repository operations needed by other packages
// type UserRepositoryInterface interface {
// 	GetUserByID(id uint, include []string) (*User, error)
// 	AssignRolesToUser(userID uint, roleIDs []uint) error
// }

// // User represents a basic user structure for shared interfaces
// type User struct {
// 	ID    uint   `json:"id"`
// 	Email string `json:"email"`
// 	Name  string `json:"name"`
// }
