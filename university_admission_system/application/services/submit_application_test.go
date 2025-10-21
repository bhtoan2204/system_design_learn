package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"university_admission_system/application/services"
	"university_admission_system/domain"
	"university_admission_system/infrastructure/repository/memory"
	appErrors "university_admission_system/pkg/errors"
	"university_admission_system/pkg/validator"
)

type (
	fixedClock struct {
		t time.Time
	}
	fixedIDGen struct {
		id string
	}
)

func (f fixedClock) Now() time.Time { return f.t }
func (f fixedIDGen) NewID() string  { return f.id }

func TestSubmitApplication_Success(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)

	applicantRepo := memory.NewApplicantRepository()
	applicationRepo := memory.NewApplicationRepository()

	applicant := &domain.Applicant{
		ID:            "applicant-1",
		FullName:      "Test User",
		Email:         "test@example.com",
		HighSchoolGPA: 3.2,
		EntranceScore: 82,
	}
	if err := applicantRepo.Save(ctx, applicant); err != nil {
		t.Fatalf("unexpected error saving applicant: %v", err)
	}

	service := services.NewSubmitApplicationService(
		applicantRepo,
		applicationRepo,
		fixedIDGen{id: "application-1"},
		fixedClock{t: now},
		validator.New(),
	)

	result, err := service.Submit(ctx, services.SubmitApplicationCommand{
		ApplicantID: applicant.ID,
		ProgramID:   "computer-science",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ApplicationID != "application-1" {
		t.Fatalf("unexpected application id %s", result.ApplicationID)
	}
	if !result.SubmittedAt.Equal(now) {
		t.Fatalf("expected submitted at %s, got %s", now, result.SubmittedAt)
	}

	stored, err := applicationRepo.FindByID(ctx, "application-1")
	if err != nil {
		t.Fatalf("unexpected error retrieving application: %v", err)
	}
	if stored == nil {
		t.Fatal("expected application to be stored")
	}
	if stored.Status != domain.ApplicationStatusSubmitted {
		t.Fatalf("expected status submitted, got %s", stored.Status)
	}
}

func TestSubmitApplication_InvalidInput(t *testing.T) {
	service := services.NewSubmitApplicationService(
		memory.NewApplicantRepository(),
		memory.NewApplicationRepository(),
		fixedIDGen{id: "application-1"},
		fixedClock{t: time.Now().UTC()},
		validator.New(),
	)

	_, err := service.Submit(context.Background(), services.SubmitApplicationCommand{})
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
	if !errors.Is(err, appErrors.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}
