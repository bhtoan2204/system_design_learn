package domain

import "time"

// EnrollmentStatus identifies whether the applicant confirmed the enrollment.
type EnrollmentStatus string

const (
	EnrollmentStatusPending   EnrollmentStatus = "pending"
	EnrollmentStatusConfirmed EnrollmentStatus = "confirmed"
)

// Enrollment captures the confirmation that the applicant will join the program.
type Enrollment struct {
	ID            string
	ApplicationID string
	OfferID       string
	Status        EnrollmentStatus
	ConfirmedAt   time.Time
}

// Confirm locks the enrollment once the applicant accepts.
func (e *Enrollment) Confirm(now time.Time) error {
	if e.Status == EnrollmentStatusConfirmed {
		return ErrEnrollmentAlreadyConfirmed
	}
	e.Status = EnrollmentStatusConfirmed
	e.ConfirmedAt = now
	return nil
}
