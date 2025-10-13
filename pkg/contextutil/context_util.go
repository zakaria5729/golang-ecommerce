package contextutil

import (
	"context"
	"errors"

	c "github.com/easy-comerce/backend/pkg/constants"
)

func GetUserIDFromContext(ctx context.Context) (*uint, error) {
	userID, ok := ctx.Value(c.UserIDContextKey).(uint)
	if !ok || userID == 0 {
		return nil, errors.New("user ID not found or invalid type in context")
	}
	return &userID, nil
}
