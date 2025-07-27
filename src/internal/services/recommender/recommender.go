package recommender

import (
	"sort"
)

type Contact struct {
	ID        int64
	Embedding []float64
}

func (c Contact) GetEmbedding() []float64 { return c.Embedding }
func (c Contact) GetID() int64            { return c.ID }

type House struct {
	ID        int64
	Embedding []float64
}

func (h House) GetEmbedding() []float64 { return h.Embedding }
func (h House) GetID() int64            { return h.ID }

type ScoredItem[T any] struct {
	Item  int64
	Score float64
}

type ScoredContact struct {
	ContactID int64
	Score     float64
}

type Embeddable interface {
	GetEmbedding() []float64
	GetID() int64
}

func RecommendContacts[S Embeddable, T Embeddable](source S, targets []T, topN int) []ScoredItem[T] {
	scored := make([]ScoredItem[T], 0, len(targets))

	for _, target := range targets {
		score := CosineSimilarity(source.GetEmbedding(), target.GetEmbedding())
		scored = append(scored, ScoredItem[T]{Item: target.GetID(), Score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if topN > len(scored) {
		topN = len(scored)
	}

	return scored[:topN]
}
