package users

import (
	"context"

	"github.com/bqdanh/money_transfer/api/grpc/user_service"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils/exceptions_parser"
	"github.com/bqdanh/money_transfer/internal/applications/users/create_user"
	"github.com/bqdanh/money_transfer/pkg/logger"
)

func (s *UserService) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	u, err := s.App.CreateUserHandler.Handle(context.Background(), create_user.CreateUserParams{
		UserName: req.GetUsername(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
		Phone:    req.GetPhone(),
	})
	if err != nil {
		logger.FromContext(ctx).Errorw("create user got exception", "err", err)
		return nil, exceptions_parser.Err2GrpcStatus(err).Err()
	}
	return &user_service.CreateUserResponse{
		Code:    utils.CodeSuccess,
		Message: utils.MessageSuccess,
		Data: &user_service.CreateUserResponse_Data{
			UserId: u.ID,
		},
	}, nil
}
