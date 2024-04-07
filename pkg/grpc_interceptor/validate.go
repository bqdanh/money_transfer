package grpc_interceptor

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The validate interface starting with protoc-gen-validate v0.6.0.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/455.
type validator interface {
	Validate(all bool) error
}

// The validate interface prior to protoc-gen-validate v0.6.0.
type validatorLegacy interface {
	Validate() error
}

type errorList interface {
	AllErrors() []error
}

type errorValidation interface {
	Field() string
	Reason() string
}

func validate(req interface{}) error {
	switch v := req.(type) {
	case validatorLegacy:
		if err := v.Validate(); err != nil {
			return parserError(err)
		}
	case validator:
		if err := v.Validate(false); err != nil {
			return parserError(err)
		}
	}
	return nil
}

func parserError(err error) error {
	if errs, ok := err.(errorList); ok {
		return parserErrorList(errs)
	}
	if v, ok := err.(errorValidation); ok {
		return parserErrorValidation(v)
	}
	return status.Error(codes.InvalidArgument, err.Error())
}

func parserErrorList(errs errorList) error {
	sts := status.New(codes.InvalidArgument, "yêu cầu không hợp lệ.")
	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0, len(errs.AllErrors()))
	for _, e := range errs.AllErrors() {
		if v, ok := e.(errorValidation); ok {
			fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       v.Field(),
				Description: v.Reason(),
			})
		}
	}
	br := &errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}
	sts, _ = sts.WithDetails(br)
	return sts.Err()
}

func parserErrorValidation(err errorValidation) error {
	sts := status.New(codes.InvalidArgument, "yêu cầu không hợp lệ.")
	br := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       err.Field(),
				Description: err.Reason(),
			},
		},
	}
	sts, _ = sts.WithDetails(br)
	return sts.Err()
}

// RequestValidationUnaryServerInterceptor returns a new unary server interceptor that validates incoming messages.
//
// Invalid messages will be rejected with `InvalidArgument` before reaching any userspace handlers.
func RequestValidationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
