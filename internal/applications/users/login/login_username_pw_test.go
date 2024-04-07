package login

import (
	"context"
	"fmt"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewLoginWithUsernamePassword(t *testing.T) {
	type args struct {
		ur userRepository
	}
	tests := []struct {
		name    string
		args    args
		want    LoginWithUsernamePassword
		wantErr error
	}{
		{
			name: "nil user repository",
			args: args{
				ur: nil,
			},
			want:    LoginWithUsernamePassword{},
			wantErr: exceptions.NewInvalidArgumentError("UserRepository", "user repository must not nil", nil),
		},
		{
			name: "valid",
			args: args{
				ur: NewMockuserRepository(gomock.NewController(t)),
			},
			want:    LoginWithUsernamePassword{userRepo: NewMockuserRepository(gomock.NewController(t))},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoginWithUsernamePassword(tt.args.ur)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHandle(t *testing.T) {
	type dependencies struct {
		userRepoFunc func(tt *testing.T) userRepository
	}
	userMock := user.User{
		ID:       1,
		UserName: "banhdanh",
		Password: "$2a$10$lepW1ozns2uifq0SNs.sduKceRysA3WvhL/puF.mMlDXMHdMprKza",
		FullName: "Banh Quoc Danh",
		Phone:    "",
	}
	mockErr := fmt.Errorf("mock error")
	tests := []struct {
		name         string
		dependencies dependencies
		args         LoginWithUsernamePasswordParams
		want         user.User
		wantErr      error
	}{
		{
			name: "login success",
			dependencies: dependencies{
				userRepoFunc: func(tt *testing.T) userRepository {
					m := NewMockuserRepository(gomock.NewController(tt))
					m.EXPECT().GetUserByUsername(gomock.Any(), "banhdanh").Return(userMock, nil)
					return m
				},
			},
			args: LoginWithUsernamePasswordParams{
				Username: "banhdanh",
				Password: "123",
			},
			want:    userMock,
			wantErr: nil,
		},
		{
			name: "repository got error",
			dependencies: dependencies{
				userRepoFunc: func(tt *testing.T) userRepository {
					m := NewMockuserRepository(gomock.NewController(tt))
					m.EXPECT().GetUserByUsername(gomock.Any(), "banhdanh").Return(user.User{}, mockErr)
					return m
				},
			},
			args: LoginWithUsernamePasswordParams{
				Username: "banhdanh",
				Password: "123",
			},
			want:    user.User{},
			wantErr: mockErr,
		},
		{
			name: "invalid password",
			dependencies: dependencies{
				userRepoFunc: func(tt *testing.T) userRepository {
					m := NewMockuserRepository(gomock.NewController(tt))
					m.EXPECT().GetUserByUsername(gomock.Any(), "banhdanh").Return(userMock, nil)
					return m
				},
			},
			args: LoginWithUsernamePasswordParams{
				Username: "banhdanh",
				Password: "xyz",
			},
			want:    user.User{},
			wantErr: exceptions.NewPreconditionError(exceptions.PreconditionTypePasswordNotMatch, exceptions.SubjectUser, "password not match", nil),
		},
		{
			name: "invalid username argument",
			dependencies: dependencies{
				userRepoFunc: func(tt *testing.T) userRepository {
					m := NewMockuserRepository(gomock.NewController(tt))
					return m
				},
			},
			args: LoginWithUsernamePasswordParams{
				Username: "",
				Password: "123",
			},
			want:    user.User{},
			wantErr: exceptions.NewInvalidArgumentError("Username", "username must not empty", nil),
		},
		{
			name: "invalid password argument",
			dependencies: dependencies{
				userRepoFunc: func(tt *testing.T) userRepository {
					m := NewMockuserRepository(gomock.NewController(tt))
					return m
				},
			},
			args: LoginWithUsernamePasswordParams{
				Username: "banhdanh",
				Password: "",
			},
			want:    user.User{},
			wantErr: exceptions.NewInvalidArgumentError("Password", "password must not empty", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := LoginWithUsernamePassword{
				userRepo: tt.dependencies.userRepoFunc(t),
			}
			got, err := h.Handle(context.TODO(), tt.args)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
