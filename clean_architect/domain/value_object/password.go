package valueobject

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedValue string
}

const (
	MinPasswordLength = 6
	MaxPasswordLength = 72 // bcrypt limit
)

func NewPassword(plainPassword string) (*Password, error) {
	if len(plainPassword) < MinPasswordLength {
		return nil, errors.New("password must be at least 6 characters")
	}
	if len(plainPassword) > MaxPasswordLength {
		return nil, errors.New("password must be at most 72 characters")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Password{hashedValue: string(hashedBytes)}, nil
}

func NewPasswordFromHash(hashedPassword string) *Password {
	return &Password{hashedValue: hashedPassword}
}

func (p *Password) Value() string {
	return p.hashedValue
}

func (p *Password) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}
