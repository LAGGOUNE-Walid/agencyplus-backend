package contact

import (
	"database/sql"
	"errors"
	"fmt"
	"logispro/internal/db"
	"logispro/internal/services/contact_service"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"logispro/internal/web/requests"
	"net/http"
	"strconv"
)

type ContactController struct {
	CreateContactService *contact_service.CreateContactService
	GetContactService    *contact_service.GetContactService
	DeleteContactService *contact_service.DeleteContactService
}

func (c *ContactController) CreateContactHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	contactService := c.CreateContactService
	req, validationErrors, err := requests.ParseCreateContactRequest(r, contactService.Queries, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}
	if len(validationErrors) > 0 {
		return response_types.ApiResponse{
			Content:    validationErrors,
			StatusCode: http.StatusBadRequest,
		}
	}

	userId, err := utils.GetUserIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	req.UserID = userId

	contactId, err := contactService.Create(r.Context(), req)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return response_types.ApiResponse{
		Content: map[string]any{
			"contact_id": contactId,
		},
		StatusCode: http.StatusCreated,
		Error:      nil,
	}
}

func (c *ContactController) GetContactsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, err := utils.GetUserIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
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
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.GetContactService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	contacts, err := c.GetContactService.All(agencyUsersId, rootId, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    contacts,
	}
}
func (c *ContactController) CountContactsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, err := utils.GetUserIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
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
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.GetContactService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	count, err := c.GetContactService.Count(agencyUsersId, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    count,
	}
}

func (c *ContactController) GetContactsListHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	userId, err := utils.GetUserIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	contactsIdsString := r.URL.Query()["ids[]"]
	if len(contactsIdsString) == 0 {
		return response_types.ApiResponse{
			StatusCode: http.StatusOK,
			Content:    nil,
		}
	}
	contactsIds := make([]int64, 0, len(contactsIdsString))
	for _, idStr := range contactsIdsString {
		idInt := utils.ParseInt64(idStr)
		contactsIds = append(contactsIds, idInt)
	}
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.GetContactService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	contacts, err := c.GetContactService.FindAll(contactsIds, agencyUsersId, r.Context())
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    contacts,
	}
}

func (c *ContactController) GetContactHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid contact ID"),
		}
	}
	userId, err := utils.GetUserIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
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
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.GetContactService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	contact, err := c.GetContactService.Get(id, agencyUsersId, r.Context())

	if err != nil {
		if err == sql.ErrNoRows {
			return response_types.ApiResponse{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("contact not found"),
			}
		}
		return response_types.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("failed to fetch contact"),
		}
	}

	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    contact,
	}
}

func (c *ContactController) DeleteContactHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err,
		}
	}
	userId, err := utils.GetUserIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	ctx := r.Context()
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(r.Context(), c.GetContactService.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})
	err = c.DeleteContactService.Delete(id, agencyUsersId, ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				StatusCode: http.StatusNotFound,
				Content:    "Contact not found",
			}
		}
		return response_types.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return response_types.ApiResponse{
		StatusCode: http.StatusOK,
		Content:    "Contact deleted successfully",
	}
}
