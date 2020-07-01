package service

import "sync"

// RatingStore is an interface to store laptop ratings
type RatingStore interface {
	// Add adds a new laptop score to the store and returns its ratings
	Add(laptopID string, score float64) (*Rating, error)
}

type Rating struct {
	Count uint32
	Sum   float64
}

// InMemoryRatingStore stores laptop ratings in memory
type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

// NewInMemoryRatingStore return a new InMemoryRatingStore
func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

// Add adds a new laptop rating to the store and returns its rating
func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopID]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}
	store.rating[laptopID] = rating
	return rating, nil
}
