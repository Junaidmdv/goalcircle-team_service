package validater

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidationError(errs validator.ValidationErrorsTranslations) (*status.Status, error) {
	st := status.New(codes.InvalidArgument, "validation error")

	var validationErrors []*errdetails.BadRequest_FieldViolation

	for field, msg := range errs {
		parts := strings.SplitN(field, ".", 2)

		if len(parts) == 2 {
			validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
				Field:       parts[1],
				Description: msg,
			})
		} else {
			validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: msg,
			})
		}

	}
	stWithDetails, err := st.WithDetails(
		&errdetails.BadRequest{
			FieldViolations: validationErrors,
		},
	)

	if err != nil {
		return nil, nil
	}

	return stWithDetails, nil
}
