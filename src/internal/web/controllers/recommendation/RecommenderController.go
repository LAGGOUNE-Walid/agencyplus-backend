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
	var params db.GetContactsWithEmbeddingsParams
	if rootId != nil {
		params = db.GetContactsWithEmbeddingsParams{
			UserID: userId,
			UserID2: sql.NullInt64{
				Valid: true,
				Int64: *rootId,
			},
		}
	} else {
		params = db.GetContactsWithEmbeddingsParams{
			UserID:  userId,
			UserID2: sql.NullInt64{Valid: false},
		}
	}
	contacts, err := c.Queries.GetContactsWithEmbeddings(ctx, params)
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
	recommendedContact := recommender.RecommendContacts(recommendedHouse, recommendingContacts, 10)

	return response_types.ApiResponse{
		Content:    recommendedContact,
		StatusCode: http.StatusOK,
	}
}
