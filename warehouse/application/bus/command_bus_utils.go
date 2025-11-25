package bus

import (
	"errors"
	"reflect"
)

func (b *bus) validate(cmd interface{}) error {
	value := reflect.ValueOf(cmd)

	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer to commands are allowed")
	}

	return nil
}
