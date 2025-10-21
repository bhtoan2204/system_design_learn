package domain

import (
	"time"
)

// ApplicationStatus captures the lifecycle of an admission application.
type ApplicationStatus string

const (
	ApplicationStatusDraft       ApplicationStatus = "draft"
	ApplicationStatusSubmitted   ApplicationStatus = "submitted"
	ApplicationStatusScored      ApplicationStatus = "scored"
	ApplicationStatusOfferIssued ApplicationStatus = "offer_issued"
)

// Application aggregates the information of a submission for admission.
type Application struct {
	ID          string
	ApplicantID string
	ProgramID   string
	Score       float64
	Status      ApplicationStatus
	SubmittedAt time.Time
	ScoredAt    time.Time
}

// Submit transitions an application from draft to submitted.
func (a *Application) Submit(submittedAt time.Time) error {
	if a.Status != ApplicationStatusDraft {
		return ErrApplicationAlreadySubmitted
	}
	a.Status = ApplicationStatusSubmitted
	a.SubmittedAt = submittedAt
	return nil
}

// RecordScore attaches the computed score to the application.
func (a *Application) RecordScore(score float64, scoredAt time.Time) error {
	switch a.Status {
	case ApplicationStatusDraft:
		return ErrApplicationNotSubmitted
	case ApplicationStatusScored, ApplicationStatusOfferIssued:
		return ErrApplicationAlreadyScored
	case ApplicationStatusSubmitted:
		// valid transition
	default:
		return ErrApplicationNotSubmitted
	}
	a.Score = score
	a.Status = ApplicationStatusScored
	a.ScoredAt = scoredAt
	return nil
}

// MarkOfferIssued marks the application as having an offer.
func (a *Application) MarkOfferIssued() error {
	if a.Status != ApplicationStatusScored {
		if a.Status == ApplicationStatusOfferIssued {
			return ErrApplicationAlreadyOffered
		}
		return ErrApplicationNotScored
	}
	a.Status = ApplicationStatusOfferIssued
	return nil
}
