package exceptions_parser

import (
	"fmt"

	"github.com/bqdanh/money_transfer/api/grpc/errdetails_custom"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fromPreconditionExceptions2GrpcStatus(perr *exceptions.PreconditionError) *status.Status {
	md := map[string]string{
		"description": perr.Description,
	}
	for k, v := range perr.Metadata {
		md[k] = fmt.Sprintf("%v", v)
	}

	ed := &errdetails_custom.FailPrecondition{
		FailureViolations: []*errdetails_custom.FailPrecondition_FailureViolation{
			{
				Reason:      string(perr.Reason),
				Subject:     perr.Subject,
				Description: perr.Description,
				Metadata:    md,
			},
		},
	}

	sts := status.New(codes.FailedPrecondition, getPrecondtionReason2Message(perr))
	sts, _ = sts.WithDetails(ed)
	return sts
}

func getPrecondtionReason2Message(perr *exceptions.PreconditionError) string {
	if msg, ok := mappingConditionToMessage[perr.Reason]; ok {
		return msg
	}
	return defautlPreconditionMessage
}

const (
	defautlPreconditionMessage = "Yêu cầu bị từ chối vì không thỏa điều kiện."
)

var mappingConditionToMessage = map[exceptions.PreconditionReason]string{
	exceptions.PreconditionReasonUserDuplicatedUserName: "Tên người dùng đã tồn tại.",
	exceptions.PreconditionReasonCannotChangeUserID:     "Không thể thay đổi ID người dùng.",
	exceptions.PreconditionReasonUserNotFound:           "Không tìm thấy người dùng.",
	exceptions.PreconditionReasonPasswordNotMatch:       "Mật khẩu không hợp lệ.",
}
