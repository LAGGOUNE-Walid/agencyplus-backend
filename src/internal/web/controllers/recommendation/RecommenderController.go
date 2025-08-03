package recommendation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"logispro/internal/db"
	"logispro/internal/services/recommender"
	"logispro/internal/shared/response_types"
	"logispro/internal/utils"
	"net/http"
	"strconv"
)

type RecommenderController struct {
	Queries *db.Queries
}

func (c *RecommenderController) GetForBuildingHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	idStr := r.PathValue("building_id")
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
	rootId, err := utils.GetRootIdFromContext(r.Context())
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	ctx := r.Context()
	buildingEmbedding, err := c.Queries.GetBuildingEmbeddings(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusOK,
			}
		}
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to get buidling embeddings id %v err %v", id, err),
			StatusCode: http.StatusInternalServerError,
		}
	}
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
	contacts, err := c.Queries.GetContactsWithEmbeddings(ctx, agencyUsersId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusOK,
			}
		}
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to get contacts with embeddings for user id %v err %v", userId, err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var recommendingContacts []recommender.Contact
	for _, contact := range contacts {
		contactId := contact.ID.Int64
		var contactEmbeddingFloat []float64
		json.Unmarshal([]byte(contact.Embedding), &contactEmbeddingFloat)
		recommendingContacts = append(recommendingContacts, recommender.Contact{ID: contactId, Embedding: contactEmbeddingFloat})
	}
	var houseEmbeddingFloat []float64
	json.Unmarshal([]byte(buildingEmbedding.Embedding), &houseEmbeddingFloat)
	recommendedHouse := recommender.House{ID: buildingEmbedding.BuildingID, Embedding: houseEmbeddingFloat}
	recommendedContacts := recommender.RecommendContacts(recommendedHouse, recommendingContacts, 10)

	return response_types.ApiResponse{
		Content:    recommendedContacts,
		StatusCode: http.StatusOK,
	}
}

func (c *RecommenderController) GetForContactsHandler(w http.ResponseWriter, r *http.Request) response_types.ApiResponse {
	ctx := r.Context()
	idStr := r.PathValue("contact_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response_types.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err,
		}
	}
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	rootId, err := utils.GetRootIdFromContext(ctx)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsers, err := utils.GetAgencyUsers(ctx, c.Queries, userId, rootId)
	if err != nil {
		return response_types.ApiResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
		return u.ID
	})

	buildings, err := c.Queries.GetBuildingsWithEmbeddings(ctx, agencyUsersId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusOK,
			}
		}
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to get buidling embeddings id %v err %v", id, err),
			StatusCode: http.StatusInternalServerError,
		}
	}
	contactEmbedding, err := c.Queries.GetContactEmbeddings(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response_types.ApiResponse{
				Content:    nil,
				StatusCode: http.StatusOK,
			}
		}
		return response_types.ApiResponse{
			Error:      fmt.Errorf("failed to get contacts with embeddings for user id %v err %v", userId, err),
			StatusCode: http.StatusInternalServerError,
		}
	}
	var recommendingBuildings []recommender.House
	for _, building := range buildings {
		buildingId := building.ID.Int64
		var buildingEmbeddingFloat []float64
		json.Unmarshal([]byte(building.Embedding), &buildingEmbeddingFloat)
		recommendingBuildings = append(recommendingBuildings, recommender.House{ID: buildingId, Embedding: buildingEmbeddingFloat})
	}

	var contactEmbeddingFloat []float64
	json.Unmarshal([]byte(contactEmbedding.Embedding), &contactEmbeddingFloat)
	recommendedContact := recommender.Contact{ID: contactEmbedding.ContactID, Embedding: contactEmbeddingFloat}
	recommendedHouses := recommender.RecommendContacts(recommendedContact, recommendingBuildings, 10)
	return response_types.ApiResponse{
		Content:    recommendedHouses,
		StatusCode: http.StatusOK,
	}
}
