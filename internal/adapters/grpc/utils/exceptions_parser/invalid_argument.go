package exceptions_parser

import (
	"fmt"

	"github.com/bqdanh/money_transfer/api/grpc/errdetails_custom"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fromInvalidArgumentExceptions2GrpcStatus(ierr *exceptions.InvalidArgumentError) *status.Status {
	md := make(map[string]string)
	for k, v := range ierr.Metadata {
		md[k] = fmt.Sprintf("%v", v)
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
	sts := status.New(codes.InvalidArgument, "Yêu cầu không hợp lệ.")
	sts, _ = sts.WithDetails(ed)
	return sts
}
