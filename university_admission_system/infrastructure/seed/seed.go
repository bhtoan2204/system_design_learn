package seed

import (
	"context"
	"time"

	"university_admission_system/domain"
	"university_admission_system/infrastructure/repository/memory"
)

// Summary captures the identifiers created by the seeding routine.
type Summary struct {
	ApplicantIDs   []string
	ApplicationIDs []string
	OfferIDs       []string
	EnrollmentIDs  []string
}

// SeedData populates the in-memory repositories with predictable demo data.
func SeedData(
	ctx context.Context,
	idGen domain.IDGenerator,
	applicants *memory.ApplicantRepository,
	applications *memory.ApplicationRepository,
	offers *memory.OfferRepository,
	enrollments *memory.EnrollmentRepository,
) (*Summary, error) {
	applicants.Reset()
	applications.Reset()
	offers.Reset()
	enrollments.Reset()

	now := time.Now().UTC()

	alice := &domain.Applicant{
		ID:            idGen.NewID(),
		FullName:      "Alice Nguyen",
		Email:         "alice@example.com",
		HighSchoolGPA: 3.6,
		EntranceScore: 85,
	}
	bob := &domain.Applicant{
		ID:            idGen.NewID(),
		FullName:      "Bob Tran",
		Email:         "bob@example.com",
		HighSchoolGPA: 2.8,
		EntranceScore: 68,
	}

	if err := applicants.Save(ctx, alice); err != nil {
		return nil, err
	}
	if err := applicants.Save(ctx, bob); err != nil {
		return nil, err
	}

	aliceApp := &domain.Application{
		ID:          idGen.NewID(),
		ApplicantID: alice.ID,
		ProgramID:   "computer-science",
		Status:      domain.ApplicationStatusDraft,
	}
	aliceSubmittedAt := now.Add(-48 * time.Hour)
	if err := aliceApp.Submit(aliceSubmittedAt); err != nil {
		return nil, err
	}
	if err := aliceApp.RecordScore(88, aliceSubmittedAt.Add(2*time.Hour)); err != nil {
		return nil, err
	}
	if err := aliceApp.MarkOfferIssued(); err != nil {
		return nil, err
	}

	bobApp := &domain.Application{
		ID:          idGen.NewID(),
		ApplicantID: bob.ID,
		ProgramID:   "business-administration",
		Status:      domain.ApplicationStatusDraft,
	}
	bobSubmittedAt := now.Add(-24 * time.Hour)
	if err := bobApp.Submit(bobSubmittedAt); err != nil {
		return nil, err
	}

	if err := applications.Save(ctx, aliceApp); err != nil {
		return nil, err
	}
	if err := applications.Save(ctx, bobApp); err != nil {
		return nil, err
	}

	alice.TrackSubmission(aliceApp.ID)
	bob.TrackSubmission(bobApp.ID)
	if err := applicants.Save(ctx, alice); err != nil {
		return nil, err
	}
	if err := applicants.Save(ctx, bob); err != nil {
		return nil, err
	}

	offerIssuedAt := now.Add(-12 * time.Hour)
	offer := &domain.Offer{
		ID:            idGen.NewID(),
		ApplicationID: aliceApp.ID,
		Score:         aliceApp.Score,
		Status:        domain.OfferStatusPending,
		IssuedAt:      offerIssuedAt,
		ExpiresAt:     now.Add(5 * 24 * time.Hour),
	}
	if err := offer.Accept(offerIssuedAt.Add(2 * time.Hour)); err != nil {
		return nil, err
	}
	if err := offers.Save(ctx, offer); err != nil {
		return nil, err
	}

	enrollment := &domain.Enrollment{
		ID:            idGen.NewID(),
		ApplicationID: aliceApp.ID,
		OfferID:       offer.ID,
		Status:        domain.EnrollmentStatusPending,
	}
	if err := enrollment.Confirm(now.Add(-6 * time.Hour)); err != nil {
		return nil, err
	}
	if err := enrollments.Save(ctx, enrollment); err != nil {
		return nil, err
	}

	return &Summary{
		ApplicantIDs:   []string{alice.ID, bob.ID},
		ApplicationIDs: []string{aliceApp.ID, bobApp.ID},
		OfferIDs:       []string{offer.ID},
		EnrollmentIDs:  []string{enrollment.ID},
	}, nil
}
