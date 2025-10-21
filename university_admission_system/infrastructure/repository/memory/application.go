package memory

import (
	"context"
	"sync"

	"university_admission_system/domain"
)

// ApplicationRepository stores applications in memory.
type ApplicationRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Application
}

// NewApplicationRepository constructs a new repository.
func NewApplicationRepository() *ApplicationRepository {
	return &ApplicationRepository{
		items: make(map[string]*domain.Application),
	}
}

// Reset clears all applications.
func (r *ApplicationRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*domain.Application)
}

// FindByID fetches an application by id.
func (r *ApplicationRepository) FindByID(_ context.Context, id string) (*domain.Application, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	app, ok := r.items[id]
	if !ok {
		return nil, nil
	}
	copy := *app
	return &copy, nil
}

// ListAll returns every application in storage.
func (r *ApplicationRepository) ListAll(_ context.Context) ([]*domain.Application, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Application, 0, len(r.items))
	for _, app := range r.items {
		copy := *app
		result = append(result, &copy)
	}
	return result, nil
}

// Save upserts the application.
func (r *ApplicationRepository) Save(_ context.Context, application *domain.Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *application
	r.items[application.ID] = &copy
	return nil
}
