package utils

import (
	"context"
	"database/sql"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
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

func GetAgencyUsers(ctx context.Context, c *db.Queries, userId int64, rootId *int64) (users []db.User, err error) {
	if rootId == nil {
		users, err = c.GetAgencyUsers(ctx, sql.NullInt64{Valid: true, Int64: userId})
		user, err := c.GetUserById(ctx, userId)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	} else {
		users, err = c.GetAgencyUsers(ctx, sql.NullInt64{Valid: true, Int64: *rootId})
		user, err := c.GetUserById(ctx, *rootId)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err
}

func ExtractField[T any, V any](input []T, extractor func(T) V) []V {
	var result []V
	for _, item := range input {
		result = append(result, extractor(item))
	}
	return result
}
