package users

import (
	"context"

	"github.com/bqdanh/money_transfer/api/grpc/user_service"
	"github.com/bqdanh/money_transfer/internal/applications/users/create_user"
	"github.com/bqdanh/money_transfer/internal/ports/grpc/common"
	"github.com/bqdanh/money_transfer/internal/ports/grpc/common/exceptions_parser"
	"github.com/bqdanh/money_transfer/pkg/logger"
)

type UserService struct {
	user_service.UnimplementedUserServiceServer
	userApp userApplications
}

type userApplications struct {
	createUser create_user.CreateUser
}

func (s *UserService) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	u, err := s.userApp.createUser.Handle(context.Background(), create_user.CreateUserParams{
		UserName: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
	})
	if err != nil {
		logger.FromContext(ctx).Errorw("create user got exception", "err", err)
		return nil, exceptions_parser.Err2GrpcStatus(err).Err()
	}
	return &user_service.CreateUserResponse{
		Code:    common.CodeSuccess,
		Message: common.MessageSuccess,
		Data: &user_service.CreateUserResponse_Data{
			UserId: u.ID,
		},
	}, nil
}
