package user_services

import (
	"context"
	"database/sql"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
)

type CreateUserService struct {
	Queries *db.Queries
}

func (s *CreateUserService) Create(req requests.CreateUserRequest) (int64, string, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, "", err
	}

	// (Optional) Save logo to disk or cloud and store filename
	var logoPath string
	if req.AgencyLogo != nil && req.LogoHeader != nil {
		logoPath, err = utils.SaveFile(req.AgencyLogo, req.LogoHeader, "uploads", constants.AgencyLogoMaxSize)
		if err != nil {
			return 0, "", err
		}
	}

	// Insert into DB using sqlc
	arg := db.CreateUserParams{
		Role:          constants.ROLE_OWENER,
		Fullname:      req.FullName,
		Email:         req.Email,
		Phone:         req.Phone,
		AgencyName:    req.AgencyName,
		AgencyAddress: req.AgencyAddress,
		AgencyLogo:    sql.NullString{String: logoPath, Valid: logoPath != ""},
		Wilaya:        req.Wilaya,
		Daira:         req.Daira,
		Password:      hashedPassword,
		RootID: func() sql.NullInt64 {
			if req.RootId != nil {
				return sql.NullInt64{Int64: *req.RootId, Valid: true}
			}
			return sql.NullInt64{Valid: false}
		}(),
	}
	res, err := s.Queries.CreateUser(context.Background(), arg)
	if err != nil {
		return 0, "", err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	token, err := utils.GenerateJWT(lastId, arg.RootID, arg.Email, arg.Role)
	if err != nil {
		return 0, "", err
	}

	return lastId, token, nil
}
