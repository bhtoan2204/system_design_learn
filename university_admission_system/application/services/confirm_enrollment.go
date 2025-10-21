package services

import (
	"context"
	"fmt"
	"time"

	"university_admission_system/domain"
	"university_admission_system/pkg/clock"
	appErrors "university_admission_system/pkg/errors"
	"university_admission_system/pkg/validator"
)

// ConfirmEnrollmentCommand holds identifiers necessary to confirm enrollment.
type ConfirmEnrollmentCommand struct {
	ApplicationID string `json:"applicationId" validate:"required"`
	OfferID       string `json:"offerId" validate:"required"`
}

// ConfirmEnrollmentResult reports confirmation timestamp.
type ConfirmEnrollmentResult struct {
	EnrollmentID string
	ConfirmedAt  time.Time
}

// ConfirmEnrollmentService turns an accepted offer into a confirmed enrollment.
type ConfirmEnrollmentService struct {
	applications domain.ApplicationRepository
	offers       domain.OfferRepository
	enrollments  domain.EnrollmentRepository
	idGen        domain.IDGenerator
	clock        clock.Clock
	validator    validator.Validator
}

// Confirm finalises the enrollment for an accepted offer.
func (s ConfirmEnrollmentService) Confirm(ctx context.Context, cmd ConfirmEnrollmentCommand) (*ConfirmEnrollmentResult, error) {
	if err := s.validator.Validate(cmd); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrInvalidInput, err)
	}

	application, err := s.applications.FindByID(ctx, cmd.ApplicationID)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, fmt.Errorf("%w: application %s", appErrors.ErrNotFound, cmd.ApplicationID)
	}

	offer, err := s.offers.FindByID(ctx, cmd.OfferID)
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, fmt.Errorf("%w: offer %s", appErrors.ErrNotFound, cmd.OfferID)
	}

	if offer.ApplicationID != application.ID {
		return nil, fmt.Errorf("%w: offer does not belong to application", appErrors.ErrConflict)
	}
	if offer.Status != domain.OfferStatusAccepted {
		return nil, fmt.Errorf("%w: offer not accepted", appErrors.ErrConflict)
	}

	enrollment := &domain.Enrollment{
		ID:            s.idGen.NewID(),
		ApplicationID: application.ID,
		OfferID:       offer.ID,
		Status:        domain.EnrollmentStatusPending,
	}

	now := s.clock.Now()
	if err := enrollment.Confirm(now); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrConflict, err)
	}

	if err := s.enrollments.Save(ctx, enrollment); err != nil {
		return nil, err
	}

	return &ConfirmEnrollmentResult{
		EnrollmentID: enrollment.ID,
		ConfirmedAt:  now,
	}, nil
}

// NewConfirmEnrollmentService builds the service.
func NewConfirmEnrollmentService(
	applications domain.ApplicationRepository,
	offers domain.OfferRepository,
	enrollments domain.EnrollmentRepository,
	idGen domain.IDGenerator,
	clock clock.Clock,
	validator validator.Validator,
) *ConfirmEnrollmentService {
	return &ConfirmEnrollmentService{
		applications: applications,
		offers:       offers,
		enrollments:  enrollments,
		idGen:        idGen,
		clock:        clock,
		validator:    validator,
	}
}
