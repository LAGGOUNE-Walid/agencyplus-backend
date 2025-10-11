package document

import (
	"database/sql"
	"errors"
	"fmt"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/services/document_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
	"net/http"
	"strconv"
)

type DocumentController struct {
	CreateDocumentService *document_service.CreateDocumentService
}

func (c *DocumentController) CreateDocumentHandler(w http.ResponseWriter, r *http.Request) response_types.Responder {
	ctx := r.Context()
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	req, validationErrors, err := requests.ParseUpdateBuildingDocumentsRequest(r, userId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}
	err = c.CreateDocumentService.Create(ctx, req)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusNoContent,
	}
}

func (c *DocumentController) GetDocumentsHandler(w http.ResponseWriter, r *http.Request) response_types.Responder {
	ctx := r.Context()
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.CreateDocumentService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	documents, err := c.CreateDocumentService.Queries.GetUserDocuments(ctx, agencyUsersId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}
	return response_types.ApiResponse{
		Content:    documents,
		StatusCode: http.StatusOK,
	}
}

func (c *DocumentController) DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) response_types.Responder {
	ctx := r.Context()
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to format user id %v to int64", ctx.Value(constants.UserIDContextKey)),
			StatusCode: http.StatusInternalServerError,
		}
	}
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.CreateDocumentService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	documentId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid document ID"),
		}
	}
	err = c.CreateDocumentService.Queries.DeleteUserDocument(ctx, db.DeleteUserDocumentParams{UsersID: agencyUsersId, ID: documentId})
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content:    nil,
		StatusCode: http.StatusOK,
	}
}

func (c *DocumentController) DownloadDocumentHandler(w http.ResponseWriter, r *http.Request) response_types.Responder {
	ctx := r.Context()
	token := r.PathValue("token")

	shareable, err := c.CreateDocumentService.Queries.GetShareableWithModel(ctx, token)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.FileResponse{
		Name:       shareable.Path.String,
		StatusCode: http.StatusOK,
	}
}
