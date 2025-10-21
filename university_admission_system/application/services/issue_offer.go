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

// IssueOfferCommand wraps the parameters needed to issue an offer.
type IssueOfferCommand struct {
	ApplicationID string `json:"applicationId" validate:"required"`
}

// IssueOfferResult exposes the offer information back to the caller.
type IssueOfferResult struct {
	OfferID   string
	Score     float64
	ExpiresAt time.Time
}

// IssueOfferService coordinates scoring and offer issuance.
type IssueOfferService struct {
	applications domain.ApplicationRepository
	applicants   domain.ApplicantRepository
	offers       domain.OfferRepository
	idGen        domain.IDGenerator
	scorer       domain.ScoreCalculator
	clock        clock.Clock
	validator    validator.Validator
	minScore     float64
}

// Issue computes the score and generates an offer when the applicant qualifies.
func (s IssueOfferService) Issue(ctx context.Context, cmd IssueOfferCommand) (*IssueOfferResult, error) {
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

	applicant, err := s.applicants.FindByID(ctx, application.ApplicantID)
	if err != nil {
		return nil, err
	}
	if applicant == nil {
		return nil, fmt.Errorf("%w: applicant %s", appErrors.ErrNotFound, application.ApplicantID)
	}

	now := s.clock.Now()
	score := s.scorer.Compute(*application, *applicant)

	if err := application.RecordScore(score, now); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrConflict, err)
	}

	if score < s.minScore {
		if err := s.applications.Save(ctx, application); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: score %.2f below threshold %.2f", appErrors.ErrConflict, score, s.minScore)
	}

	if err := application.MarkOfferIssued(); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrConflict, err)
	}

	offer := &domain.Offer{
		ID:            s.idGen.NewID(),
		ApplicationID: application.ID,
		Score:         score,
		Status:        domain.OfferStatusPending,
		IssuedAt:      now,
		ExpiresAt:     now.Add(7 * 24 * time.Hour),
	}

	if err := s.applications.Save(ctx, application); err != nil {
		return nil, err
	}

	if err := s.offers.Save(ctx, offer); err != nil {
		return nil, err
	}

	return &IssueOfferResult{
		OfferID:   offer.ID,
		Score:     score,
		ExpiresAt: offer.ExpiresAt,
	}, nil
}

// NewIssueOfferService constructs the service.
func NewIssueOfferService(
	applications domain.ApplicationRepository,
	applicants domain.ApplicantRepository,
	offers domain.OfferRepository,
	idGen domain.IDGenerator,
	scorer domain.ScoreCalculator,
	clock clock.Clock,
	validator validator.Validator,
	minimumScore float64,
) *IssueOfferService {
	return &IssueOfferService{
		applications: applications,
		applicants:   applicants,
		offers:       offers,
		idGen:        idGen,
		scorer:       scorer,
		clock:        clock,
		validator:    validator,
		minScore:     minimumScore,
	}
}
