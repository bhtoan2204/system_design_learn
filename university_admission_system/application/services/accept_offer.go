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

// AcceptOfferCommand contains user input to accept an offer.
type AcceptOfferCommand struct {
	OfferID string `json:"offerId" validate:"required"`
}

// AcceptOfferResult returns the relevant acceptance timestamp.
type AcceptOfferResult struct {
	AcceptedAt time.Time
}

// AcceptOfferService orchestrates offer acceptance.
type AcceptOfferService struct {
	offers    domain.OfferRepository
	clock     clock.Clock
	validator validator.Validator
}

// Accept triggers domain acceptance logic and persists the offer.
func (s AcceptOfferService) Accept(ctx context.Context, cmd AcceptOfferCommand) (*AcceptOfferResult, error) {
	if err := s.validator.Validate(cmd); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrInvalidInput, err)
	}

	offer, err := s.offers.FindByID(ctx, cmd.OfferID)
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, fmt.Errorf("%w: offer %s", appErrors.ErrNotFound, cmd.OfferID)
	}

	now := s.clock.Now()
	if err := offer.Accept(now); err != nil {
		return nil, fmt.Errorf("%w: %v", appErrors.ErrConflict, err)
	}

	if err := s.offers.Save(ctx, offer); err != nil {
		return nil, err
	}

	return &AcceptOfferResult{
		AcceptedAt: now,
	}, nil
}

// NewAcceptOfferService creates an AcceptOfferService.
func NewAcceptOfferService(
	offers domain.OfferRepository,
	clock clock.Clock,
	validator validator.Validator,
) *AcceptOfferService {
	return &AcceptOfferService{
		offers:    offers,
		clock:     clock,
		validator: validator,
	}
}
