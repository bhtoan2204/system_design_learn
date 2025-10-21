package domain

import "errors"

var (
	ErrApplicationAlreadySubmitted = errors.New("application already submitted")
	ErrApplicationNotSubmitted     = errors.New("application not submitted")
	ErrApplicationAlreadyScored    = errors.New("application already scored")
	ErrApplicationAlreadyOffered   = errors.New("application already has an offer")
	ErrApplicationNotScored        = errors.New("application not scored")
	ErrOfferExpired                = errors.New("offer already expired")
	ErrOfferAlreadyAccepted        = errors.New("offer already accepted")
	ErrEnrollmentAlreadyConfirmed  = errors.New("enrollment already confirmed")
)
