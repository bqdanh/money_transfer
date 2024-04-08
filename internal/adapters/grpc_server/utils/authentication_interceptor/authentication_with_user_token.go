package authentication_interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/bqdanh/money_transfer/api/grpc/errdetails_custom"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils/exceptions_parser"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/validate_user_token"
	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthenticationWithUserToken struct {
	validateUserToken         validate_user_token.ValidateUserToken
	methodsByPassAuthenticate map[string]struct{}
}

func NewAuthenticationWithUserToken(v validate_user_token.ValidateUserToken, methodsByPassAuthenticate []string) AuthenticationWithUserToken {
	m := make(map[string]struct{})
	for _, method := range methodsByPassAuthenticate {
		m[method] = struct{}{}
	}
	return AuthenticationWithUserToken{
		validateUserToken:         v,
		methodsByPassAuthenticate: m,
	}
}

func (a *AuthenticationWithUserToken) ValidateToken(ctx context.Context, userToken string) (authenticate.UserAuthenticateData, error) {
	data, err := a.validateUserToken.Handle(ctx, userToken)
	if err != nil {
		return authenticate.UserAuthenticateData{}, fmt.Errorf("validate user token: %w", err)
	}
	return data, nil
}

const (
	authenticationHeader = "authorization"
)

func (a *AuthenticationWithUserToken) UserTokenAuthenticationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := a.methodsByPassAuthenticate[info.FullMethod]; ok {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			sts := status.New(codes.PermissionDenied, "thông tin xác thực không hợp lệ.")
			sts, _ = sts.WithDetails(
				&errdetails_custom.FailPrecondition{
					FailureViolations: []*errdetails_custom.FailPrecondition_FailureViolation{
						{
							Reason:      string(exceptions.PreconditionReasonInvalidToken),
							Subject:     exceptions.SubjectAuthentication,
							Description: "not found token in incoming context",
							Metadata: map[string]string{
								"method": info.FullMethod,
							},
						},
					},
				},
			)
			return nil, sts.Err()
		}
		tokens := md.Get(authenticationHeader)
		if len(tokens) == 0 {
			sts := status.New(codes.PermissionDenied, "thông tin xác thực không hợp lệ.")
			sts, _ = sts.WithDetails(
				&errdetails_custom.FailPrecondition{
					FailureViolations: []*errdetails_custom.FailPrecondition_FailureViolation{
						{
							Reason:      string(exceptions.PreconditionReasonInvalidToken),
							Subject:     exceptions.SubjectAuthentication,
							Description: "not found token in header",
							Metadata: map[string]string{
								"method": info.FullMethod,
							},
						},
					},
				},
			)
			return nil, sts.Err()
		}
		token := strings.Replace(tokens[0], "Bearer ", "", 1)
		data, err := a.ValidateToken(ctx, token)
		if err != nil {
			sts := exceptions_parser.Err2GrpcStatus(err)
			return nil, sts.Err()
		}
		ctx = authenticate.ContextWithAuthentication(ctx, data)
		return handler(ctx, req)
	}
}
