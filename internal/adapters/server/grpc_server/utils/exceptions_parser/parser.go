package exceptions_parser

import (
	"errors"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Err2GrpcStatus(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "success")
	}

	sts, ok := exceptions2GrpcStatus(err)
	if ok {
		return sts
	}
	sts = status.New(codes.Internal, "lỗi hệ thống.")
	sts, _ = sts.WithDetails(&errdetails.ErrorInfo{
		Reason: "INTERNAL_ERROR",
		Domain: "money_transfer",
		Metadata: map[string]string{
			"type":      getErrorType(err),
			"raw_error": err.Error(),
		},
	})

	return sts
}

func getErrorType(err error) string {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return fmt.Sprintf("%T", err)
	}

	return getErrorType(u.Unwrap())
}

func exceptions2GrpcStatus(err error) (*status.Status, bool) {
	if err == nil {
		return nil, true
	}

	var perr *exceptions.PreconditionError
	if errors.As(err, &perr) {
		sts := fromPreconditionExceptions2GrpcStatus(perr)
		return sts, true
	}

	var ierr *exceptions.InvalidArgumentError
	if errors.As(err, &ierr) {
		sts := fromInvalidArgumentExceptions2GrpcStatus(ierr)
		return sts, true
	}

	return nil, false
}
