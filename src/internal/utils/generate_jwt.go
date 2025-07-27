package utils

import (
	"database/sql"
	"logispro/internal/config"
	"logispro/internal/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int64, rootId sql.NullInt64, email string, userRole int64) (string, error) {
	var rootId64p *int64
	if rootId.Valid {
		rootId64p = &rootId.Int64
	} else {
		rootId64p = nil
	}
	claims := jwt.MapClaims{
		constants.UserIDContextKey:    userID,
		constants.UserRoleContextKey:  userRole,
		constants.UserEmailContextKey: email,
		constants.UserRootContextKey:  rootId64p,
		"exp":                         time.Now().Add(24 * time.Hour * 365).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret)
}
