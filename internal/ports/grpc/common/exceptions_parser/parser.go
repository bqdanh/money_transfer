package exceptions_parser

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bqdanh/money_transfer/api/grpc/errdetails_custom"
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
	return status.New(codes.Internal, err.Error())
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

func fromPreconditionExceptions2GrpcStatus(perr *exceptions.PreconditionError) *status.Status {
	sts := status.New(codes.FailedPrecondition, perr.Error())
	md := map[string]string{
		"description": perr.Description,
	}
	for k, v := range perr.Metadata {
		vstr := fmt.Sprintf("%v", v)
		if bs, jerr := json.Marshal(v); jerr == nil {
			vstr = string(bs)
		}
		md[k] = vstr
	}
	ed := &errdetails.ErrorInfo{
		Reason:   string(perr.Type),
		Domain:   perr.Subject,
		Metadata: md,
	}

	sts, _ = sts.WithDetails(ed)
	return sts
}

func fromInvalidArgumentExceptions2GrpcStatus(ierr *exceptions.InvalidArgumentError) *status.Status {
	sts := status.New(codes.InvalidArgument, ierr.Error())
	md := make(map[string]string)
	for k, v := range ierr.Metadata {
		vstr := fmt.Sprintf("%v", v)
		if bs, jerr := json.Marshal(v); jerr == nil {
			vstr = string(bs)
		}
		md[k] = vstr
	}
	ed := &errdetails_custom.BadRequest{
		FieldViolations: []*errdetails_custom.BadRequest_FieldViolation{
			{
				Field:       ierr.Field,
				Description: ierr.Description,
				Metadata:    md,
			},
		},
	}

	sts, _ = sts.WithDetails(ed)
	return sts
}
