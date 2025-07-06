package recommender

import (
	"sort"
)

type Contact struct {
	ID        int64
	Embedding []float64
}

type House struct {
	ID        int64
	Embedding []float64
}

type ScoredContact struct {
	ContactID int64
	Score     float64
}

func RecommendContacts(house House, contacts []Contact, topN int) []ScoredContact {
	scored := []ScoredContact{}

	for _, contact := range contacts {
		sim := CosineSimilarity(house.Embedding, contact.Embedding)
		scored = append(scored, ScoredContact{ContactID: contact.ID, Score: sim})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if topN > len(scored) {
		topN = len(scored)
	}

	return scored[:topN]
}
