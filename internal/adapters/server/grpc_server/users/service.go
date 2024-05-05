package users

import (
	"github.com/bqdanh/money_transfer/api/grpc/user_service"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/login"
	"github.com/bqdanh/money_transfer/internal/applications/users/create_user"
	"google.golang.org/grpc"
)

type UserService struct {
	user_service.UnimplementedUserServiceServer
	App UserApplications
}

type UserApplications struct {
	CreateUserHandler create_user.CreateUser
	Login             login.Login
}

func NewUserService(app UserApplications) *UserService {
	return &UserService{
		App: app,
	}
}

func (s *UserService) RegisterService(server grpc.ServiceRegistrar) {
	user_service.RegisterUserServiceServer(server, s)
}
