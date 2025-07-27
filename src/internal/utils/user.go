package utils

import (
	"context"
	"fmt"
	"logispro/internal/constants"
)

func GetUserIdFromContext(ctx context.Context) (int64, error) {
	userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
	if !ok {
		return 0, fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey))
	}
	return userId, nil
}

func GetRootIdFromContext(ctx context.Context) (rootId *int64, err error) {
	if ctx.Value(constants.UserRootContextKey) == nil {
		return rootId, nil
	}
	rootId, ok := ctx.Value(constants.UserRootContextKey).(*int64)
	if !ok {
		return rootId, fmt.Errorf("failed to format root id %v to *int64", ctx.Value(constants.UserRootContextKey))
	}
	return rootId, nil
}
