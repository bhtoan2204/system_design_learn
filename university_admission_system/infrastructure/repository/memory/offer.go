package memory

import (
	"context"
	"sync"

	"university_admission_system/domain"
)

// OfferRepository persists offers using a map.
type OfferRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Offer
}

// NewOfferRepository builds a new instance.
func NewOfferRepository() *OfferRepository {
	return &OfferRepository{
		items: make(map[string]*domain.Offer),
	}
}

// Reset wipes the repository content.
func (r *OfferRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*domain.Offer)
}

// FindByID returns the offer for the identifier.
func (r *OfferRepository) FindByID(_ context.Context, id string) (*domain.Offer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	offer, ok := r.items[id]
	if !ok {
		return nil, nil
	}
	copy := *offer
	return &copy, nil
}

// ListAll exposes every offer stored.
func (r *OfferRepository) ListAll(_ context.Context) ([]*domain.Offer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Offer, 0, len(r.items))
	for _, offer := range r.items {
		copy := *offer
		result = append(result, &copy)
	}
	return result, nil
}

// Save adds or updates an offer.
func (r *OfferRepository) Save(_ context.Context, offer *domain.Offer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *offer
	r.items[offer.ID] = &copy
	return nil
}
