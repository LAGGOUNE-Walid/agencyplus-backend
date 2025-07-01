package utils

import (
	"logispro/internal/config"
	"logispro/internal/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int64, email string, userRole int64) (string, error) {
	claims := jwt.MapClaims{
		constants.UserIDContextKey:    userID,
		constants.UserRoleContextKey:  userRole,
		constants.UserEmailContextKey: email,
		"exp":                         time.Now().Add(24 * time.Hour * 365).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret)
}
