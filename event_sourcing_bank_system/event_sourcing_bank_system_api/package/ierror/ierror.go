package ierror

import (
	"errors"
	"reflect"

	pkgErr "github.com/pkg/errors"
)

type InternalError struct {
	RootErr  error
	Msg      string
	Code     string
	HttpCode int
	GrpcCode int
}

func (e *InternalError) Error() string {
	return e.Msg
}

var (
	ErrRecordNotFound    = errorWithStack(errors.New("record NOT FOUND"))
	ErrOptimisticLock    = errorWithStack(errors.New("optimistic_lock_msg_err"))
	ErrNotHavePermission = errorWithStack(errors.New("not_have_permission_msg_err"))
	ErrUnsupported       = errorWithStack(errors.ErrUnsupported)

	ErrFieldRequired = func(field string) error {
		return errorWithStack(errors.New(field + " is required"))
	}

	ErrMissingParam = func(field string) error {
		return errorWithStack(errors.New(field + " is missing"))
	}

	ErrInvalidParam = func(field string) error {
		return errorWithStack(errors.New(field + " is invalid"))
	}
)

func errorWithStack(err error) error {
	if err == nil {
		return nil
	}

	stackTrace := reflect.ValueOf(err).MethodByName("StackTrace")
	if stackTrace.IsValid() {
		return err
	}

	return pkgErr.WithStack(err)
}

func CustomError(xErr InternalError) *InternalError {
	return &InternalError{
		RootErr:  errorWithStack(xErr.RootErr),
		Msg:      xErr.Msg,
		HttpCode: xErr.HttpCode,
		GrpcCode: xErr.GrpcCode,
	}
}

func Error(err error) error {
	if err == nil {
		return nil
	}

	return errorWithStack(err)
}
