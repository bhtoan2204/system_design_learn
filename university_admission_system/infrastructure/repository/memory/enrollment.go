package memory

import (
	"context"
	"sync"

	"university_admission_system/domain"
)

// EnrollmentRepository holds enrollments in memory.
type EnrollmentRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Enrollment
}

// NewEnrollmentRepository creates a fresh repository.
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		items: make(map[string]*domain.Enrollment),
	}
}

// Reset clears stored enrollments.
func (r *EnrollmentRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*domain.Enrollment)
}

// FindByID obtains an enrollment by identifier.
func (r *EnrollmentRepository) FindByID(_ context.Context, id string) (*domain.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, ok := r.items[id]
	if !ok {
		return nil, nil
	}
	copy := *enrollment
	return &copy, nil
}

// ListAll returns every enrollment.
func (r *EnrollmentRepository) ListAll(_ context.Context) ([]*domain.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Enrollment, 0, len(r.items))
	for _, enrollment := range r.items {
		copy := *enrollment
		result = append(result, &copy)
	}
	return result, nil
}

// Save persists the enrollment.
func (r *EnrollmentRepository) Save(_ context.Context, enrollment *domain.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *enrollment
	r.items[enrollment.ID] = &copy
	return nil
}
