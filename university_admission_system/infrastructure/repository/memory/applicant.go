package memory

import (
	"context"
	"sync"

	"university_admission_system/domain"
)

// ApplicantRepository stores applicants in a thread-safe map.
type ApplicantRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Applicant
}

// NewApplicantRepository initializes an empty repository.
func NewApplicantRepository() *ApplicantRepository {
	return &ApplicantRepository{
		items: make(map[string]*domain.Applicant),
	}
}

// Reset clears the stored applicants.
func (r *ApplicantRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*domain.Applicant)
}

// FindByID returns an applicant or nil when it cannot be found.
func (r *ApplicantRepository) FindByID(_ context.Context, id string) (*domain.Applicant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	applicant, ok := r.items[id]
	if !ok {
		return nil, nil
	}
	copy := *applicant
	return &copy, nil
}

// ListAll enumerates every applicant in storage.
func (r *ApplicantRepository) ListAll(_ context.Context) ([]*domain.Applicant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Applicant, 0, len(r.items))
	for _, applicant := range r.items {
		copy := *applicant
		result = append(result, &copy)
	}
	return result, nil
}

// Save upserts an applicant in memory.
func (r *ApplicantRepository) Save(_ context.Context, applicant *domain.Applicant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *applicant
	r.items[applicant.ID] = &copy
	return nil
}
