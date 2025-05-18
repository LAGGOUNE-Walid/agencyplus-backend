package user_services

import (
	"context"
	"database/sql"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
)

type UpdateUserService struct {
	Queries *db.Queries
}

func (s *UpdateUserService) Update(ctx context.Context, req requests.UpdateUserRequest) error {
	existingByEmail, err := s.Queries.CountUsersByEmailExcludingID(ctx, db.CountUsersByEmailExcludingIDParams{
		Email: req.Email,
		ID:    req.UserID,
	})
	if err != nil {
		return err
	}
	if existingByEmail > 0 {
		return fmt.Errorf("email is already used")
	}

	if req.HasPassword {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		err = s.Queries.UpdatePassword(ctx, db.UpdatePasswordParams{
			ID:       req.UserID,
			Password: hashedPassword,
		})
		if err != nil {
			return err
		}
	}

	if req.HasLogo {
		logoPath, err := utils.SaveImageFile(req.AgencyLogo, req.LogoHeader, "uploads", constants.AgencyLogoMaxSize)
		if err != nil {
			return err
		}
		err = s.Queries.UpdateLogo(ctx, db.UpdateLogoParams{
			ID:         req.UserID,
			AgencyLogo: sql.NullString{String: logoPath, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	err = s.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:            req.UserID,
		Fullname:      req.FullName,
		Email:         req.Email,
		Phone:         req.Phone,
		AgencyName:    req.AgencyName,
		AgencyAddress: req.AgencyAddress,
		Wilaya:        req.Wilaya,
		Daira:         req.Daira,
	})
	if err != nil {
		return err
	}

	return nil
}
