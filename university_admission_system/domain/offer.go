package domain

import "time"

// OfferStatus describes whether the offer has been accepted or expired.
type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "pending"
	OfferStatusAccepted OfferStatus = "accepted"
	OfferStatusExpired  OfferStatus = "expired"
)

// Offer represents the admission decision sent to an applicant.
type Offer struct {
	ID            string
	ApplicationID string
	Score         float64
	Status        OfferStatus
	IssuedAt      time.Time
	ExpiresAt     time.Time
	AcceptedAt    time.Time
}

// Accept marks the offer as accepted if it is still valid.
func (o *Offer) Accept(now time.Time) error {
	if o.Status == OfferStatusAccepted {
		return ErrOfferAlreadyAccepted
	}
	if now.After(o.ExpiresAt) {
		o.Status = OfferStatusExpired
		return ErrOfferExpired
	}
	o.Status = OfferStatusAccepted
	o.AcceptedAt = now
	return nil
}
