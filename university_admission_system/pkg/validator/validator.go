package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Validator exposes request validation.
type Validator interface {
	Validate(target interface{}) error
}

// SimpleValidator inspects struct tags and validates basic constraints.
type SimpleValidator struct{}

// New creates a new SimpleValidator instance.
func New() *SimpleValidator {
	return &SimpleValidator{}
}

// Validate enforces `validate:"required"` on string fields.
func (SimpleValidator) Validate(target interface{}) error {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if !value.IsValid() {
		return errors.New("validator: invalid value")
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	var violations []string
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("validate")
		if !strings.Contains(tag, "required") {
			continue
		}
		fieldValue := value.Field(i)
		if fieldValue.Kind() == reflect.String && fieldValue.Len() == 0 {
			name := field.Name
			if jsonName := field.Tag.Get("json"); jsonName != "" {
				name = strings.Split(jsonName, ",")[0]
			}
			violations = append(violations, fmt.Sprintf("%s is required", name))
			continue
		}
	}

	if len(violations) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(violations, ", "))
	}
	return nil
}
