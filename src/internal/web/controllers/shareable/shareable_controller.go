package shareable

import (
	"database/sql"
	"errors"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/interfaces"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type ShareableController struct {
	Queries *db.Queries
}

type DocumentShareable struct {
	ID   int64
	path string
}

func (d DocumentShareable) GetID() int64     { return d.ID }
func (d DocumentShareable) GetType() string  { return "document" }
func (d DocumentShareable) GetTitle() string { return d.path }

type BuildingShareable struct {
	ID   int64
	name string
}

func (b BuildingShareable) GetID() int64     { return b.ID }
func (b BuildingShareable) GetType() string  { return "building" }
func (b BuildingShareable) GetTitle() string { return b.name }

func (c *ShareableController) Share(w http.ResponseWriter, r *http.Request) response_types.Responder {
	ctx := r.Context()
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid building ID"),
		}
	}
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	shareableType := r.PathValue("type")
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	var shareable interfaces.Shareable
	if shareableType == "document" {
		document, err := c.Queries.GetDocumentById(ctx, db.GetDocumentByIdParams{ID: id, UsersID: agencyUsersId})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return response_types.ApiResponse{
					StatusCode: http.StatusNotFound,
					Error:      fmt.Errorf("document not defined"),
				}
			}
			return response_types.ApiResponse{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}
		shareable = DocumentShareable{
			ID:   document.ID,
			path: document.Path,
		}
	} else if shareableType == "building" {
		building, err := c.Queries.GetBuilding(ctx, db.GetBuildingParams{ID: id, UsersID: agencyUsersId})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return response_types.ApiResponse{
					StatusCode: http.StatusNotFound,
					Error:      fmt.Errorf("building not defined"),
				}
			}
			return response_types.ApiResponse{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}
		shareable = BuildingShareable{
			ID:   building.ID,
			name: building.Title.String,
		}
	} else {
		return response_types.ApiResponse{
			StatusCode: http.StatusNotFound,
			Error:      fmt.Errorf("share type not defined"),
		}
	}
	token := uuid.New().String()
	created, err := c.Queries.CreateShareable(ctx, db.CreateShareableParams{
		Token:     token,
		ModelType: shareable.GetType(),
		ModelID:   shareable.GetID(),
		UserID:    userId,
	})
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    created,
		StatusCode: http.StatusCreated,
	}
}
