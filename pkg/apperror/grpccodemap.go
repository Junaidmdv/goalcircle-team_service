package apperror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errTypeToGRPCCode = map[ErrorType]codes.Code{
	ErrTypeUnknown:          codes.Unknown,
	ErrorTypeNotFound:       codes.NotFound,
	ErrorTypeInternal:       codes.Internal,
	ErrorTypeConflict:       codes.AlreadyExists,
	ErrorInvalidArguement:   codes.InvalidArgument,
	ErrorDeadlineExceed:     codes.DeadlineExceeded,
	ErrorUnAuthenticated:    codes.Unauthenticated,
	ErrorFailedPrecondition: codes.FailedPrecondition,
	ErrorPermissionDenied:   codes.PermissionDenied,
	ErrorBadRequest:         codes.InvalidArgument,
}

func GRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	domainErr, ok := err.(*DomainError)
	if !ok {
		return status.Error(codes.Internal, err.Error())
	}

	code, exists := errTypeToGRPCCode[domainErr.Type]
	if !exists {
		code = codes.Unknown
	}

	return status.Error(code, domainErr.Message)
}
