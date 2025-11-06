package valueobject

import (
	"errors"
	"regexp"
)

type Email struct {
	value string
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func NewEmail(value string) (*Email, error) {
	if !isValidEmail(value) {
		return nil, errors.New("invalid email")
	}
	return &Email{value: value}, nil
}

func (e *Email) Value() string {
	return e.value
}
