package user_services

import (
	"context"
	"fmt"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Queries *db.Queries
}

func (s *AuthService) Authenticate(ctx context.Context, req requests.AuthRequest) (db.User, string, error) {
	user, err := s.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return db.User{}, "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return db.User{}, "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.RootID, user.Email, user.Role)
	if err != nil {
		return db.User{}, "", fmt.Errorf("could not generate token")
	}

	return user, token, nil
}
