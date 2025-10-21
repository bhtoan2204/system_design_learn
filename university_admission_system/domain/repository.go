package domain

import "context"

// ApplicantRepository provides access to applicant aggregates.
type ApplicantRepository interface {
	FindByID(ctx context.Context, id string) (*Applicant, error)
	ListAll(ctx context.Context) ([]*Applicant, error)
	Save(ctx context.Context, applicant *Applicant) error
}

// ApplicationRepository stores admission applications.
type ApplicationRepository interface {
	FindByID(ctx context.Context, id string) (*Application, error)
	ListAll(ctx context.Context) ([]*Application, error)
	Save(ctx context.Context, application *Application) error
}

// OfferRepository manages offers.
type OfferRepository interface {
	FindByID(ctx context.Context, id string) (*Offer, error)
	ListAll(ctx context.Context) ([]*Offer, error)
	Save(ctx context.Context, offer *Offer) error
}

// EnrollmentRepository manages enrollments.
type EnrollmentRepository interface {
	FindByID(ctx context.Context, id string) (*Enrollment, error)
	ListAll(ctx context.Context) ([]*Enrollment, error)
	Save(ctx context.Context, enrollment *Enrollment) error
}
