package apperror

import (
	"errors"

	"google.golang.org/grpc/codes"
)

type ErrorType struct {
	t string
}

var (
	ErrTypeUnknown        = ErrorType{"unknown"}
	ErrTypeNotFound       = ErrorType{"not-found"}
	ErrTypeAlreadyExists  = ErrorType{"already-exists"}
	ErrTypeIncorrectInput = ErrorType{"incorrect-input"}
)

type AppError struct {
	error          error
	errorType      ErrorType
	grpcStatusCode codes.Code
}

func (ae AppError) Error() string {
	return ae.error.Error()
}
func (ae AppError) Type() ErrorType {
	return ae.errorType
}
func (ae AppError) GRPCStatusCode() codes.Code {
	return ae.grpcStatusCode
}
func (ae AppError) Is(target error) bool {
	var appError *AppError
	if errors.As(target, &appError) {
		return ae.Type().t == appError.Type().t
	}

	return false
}

func NewApplicationError(error error) AppError {
	return AppError{
		error:          error,
		errorType:      ErrTypeUnknown,
		grpcStatusCode: codes.Unknown,
	}
}
func NewNotFoundError(error error) AppError {
	return AppError{
		error:          error,
		errorType:      ErrTypeNotFound,
		grpcStatusCode: codes.NotFound,
	}
}
func NewAlreadyExistsError(error error) AppError {
	return AppError{
		error:          error,
		errorType:      ErrTypeAlreadyExists,
		grpcStatusCode: codes.AlreadyExists,
	}
}
func NewIncorrectInputError(error error) AppError {
	return AppError{
		error:          error,
		errorType:      ErrTypeIncorrectInput,
		grpcStatusCode: codes.InvalidArgument,
	}
}
