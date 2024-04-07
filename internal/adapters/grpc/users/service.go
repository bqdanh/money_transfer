package users

import (
	"context"

	"github.com/bqdanh/money_transfer/api/grpc/user_service"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc/utils"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc/utils/exceptions_parser"
	"github.com/bqdanh/money_transfer/internal/applications/users/create_user"
	"github.com/bqdanh/money_transfer/pkg/logger"
	"google.golang.org/grpc"
)

type UserService struct {
	user_service.UnimplementedUserServiceServer
	App UserApplications
}

type UserApplications struct {
	CreateUser create_user.CreateUser
}

func NewUserService(app UserApplications) *UserService {
	return &UserService{
		App: app,
	}
}

func (s *UserService) RegisterService(server grpc.ServiceRegistrar) {
	user_service.RegisterUserServiceServer(server, s)
}

func (s *UserService) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	u, err := s.App.CreateUser.Handle(context.Background(), create_user.CreateUserParams{
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
		Code:    utils.CodeSuccess,
		Message: utils.MessageSuccess,
		Data: &user_service.CreateUserResponse_Data{
			UserId: u.ID,
		},
	}, nil
}
