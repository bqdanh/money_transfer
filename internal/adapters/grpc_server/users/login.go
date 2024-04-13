package users

import (
	"context"

	"github.com/bqdanh/money_transfer/api/grpc/user_service"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils/exceptions_parser"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/login"
	"github.com/bqdanh/money_transfer/pkg/logger"
)

func (s *UserService) Login(ctx context.Context, req *user_service.LoginRequest) (*user_service.LoginResponse, error) {
	result, err := s.App.Login.Handle(ctx, login.LoginParams{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		logger.FromContext(ctx).Errorw("login got exception", "err", err)
		return nil, exceptions_parser.Err2GrpcStatus(err).Err()
	}
	return &user_service.LoginResponse{
		Code:    utils.CodeSuccess,
		Message: utils.MessageSuccess,
		Data: &user_service.LoginResponse_Data{
			Token:  result.Token,
			UserId: result.User.ID,
		},
	}, nil
}
