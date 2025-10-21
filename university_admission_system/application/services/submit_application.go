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

// SubmitApplicationCommand represents the payload for submitting an application.
type SubmitApplicationCommand struct {
	ApplicantID string `json:"applicantId" validate:"required"`
	ProgramID   string `json:"programId" validate:"required"`
}

// SubmitApplicationResult returns the newly created application identifier.
type SubmitApplicationResult struct {
	ApplicationID string
	SubmittedAt   time.Time
}

// SubmitApplicationService orchestrates the submission use-case.
type SubmitApplicationService struct {
	applicants domain.ApplicantRepository
	apps       domain.ApplicationRepository
	idGen      domain.IDGenerator
	clock      clock.Clock
	validator  validator.Validator
}

// Submit orchestrates validation and persistence of a new application.
func (s SubmitApplicationService) Submit(ctx context.Context, cmd SubmitApplicationCommand) (*SubmitApplicationResult, error) {
	if err := s.validator.Validate(cmd); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrInvalidInput, err)
	}

	applicant, err := s.applicants.FindByID(ctx, cmd.ApplicantID)
	if err != nil {
		return nil, err
	}
	if applicant == nil {
		return nil, fmt.Errorf("%w: applicant %s", appErrors.ErrNotFound, cmd.ApplicantID)
	}

	if !applicant.CanSubmit() {
		return nil, fmt.Errorf("%w: applicant not eligible", appErrors.ErrConflict)
	}

	now := s.clock.Now()
	newApplication := &domain.Application{
		ID:          s.idGen.NewID(),
		ApplicantID: cmd.ApplicantID,
		ProgramID:   cmd.ProgramID,
		Status:      domain.ApplicationStatusDraft,
	}

	if err := newApplication.Submit(now); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrConflict, err)
	}

	if err := s.apps.Save(ctx, newApplication); err != nil {
		return nil, err
	}

	applicant.TrackSubmission(newApplication.ID)
	if err := s.applicants.Save(ctx, applicant); err != nil {
		return nil, err
	}

	return &SubmitApplicationResult{
		ApplicationID: newApplication.ID,
		SubmittedAt:   now,
	}, nil
}

// NewSubmitApplicationService constructs the service with its dependencies.
func NewSubmitApplicationService(
	applicants domain.ApplicantRepository,
	applications domain.ApplicationRepository,
	idGen domain.IDGenerator,
	clock clock.Clock,
	validator validator.Validator,
) *SubmitApplicationService {
	return &SubmitApplicationService{
		applicants: applicants,
		apps:       applications,
		idGen:      idGen,
		clock:      clock,
		validator:  validator,
	}
}
