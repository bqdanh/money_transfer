package validate_username_password

import (
	"context"
	"fmt"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewValidateUsernamePassword(t *testing.T) {
	type args struct {
		ur userRepository
	}
	tests := []struct {
		name    string
		args    args
		want    ValidateUsernamePassword
		wantErr error
	}{
		{
			name: "nil user repository",
			args: args{
				ur: nil,
			},
			want:    ValidateUsernamePassword{},
			wantErr: exceptions.NewInvalidArgumentError("UserRepository", "user repository must not nil", nil),
		},
		{
			name: "valid",
			args: args{
				ur: NewMockuserRepository(gomock.NewController(t)),
			},
			want:    ValidateUsernamePassword{userRepo: NewMockuserRepository(gomock.NewController(t))},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewValidateUsernamePassword(tt.args.ur)
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
		args         ValidateUsernamePasswordParams
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
			args: ValidateUsernamePasswordParams{
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
			args: ValidateUsernamePasswordParams{
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
			args: ValidateUsernamePasswordParams{
				Username: "banhdanh",
				Password: "xyz",
			},
			want:    user.User{},
			wantErr: exceptions.NewPreconditionError(exceptions.PreconditionReasonPasswordNotMatch, exceptions.SubjectUser, "password not match", nil),
		},
		{
			name: "invalid username argument",
			dependencies: dependencies{
				userRepoFunc: func(tt *testing.T) userRepository {
					m := NewMockuserRepository(gomock.NewController(tt))
					return m
				},
			},
			args: ValidateUsernamePasswordParams{
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
			args: ValidateUsernamePasswordParams{
				Username: "banhdanh",
				Password: "",
			},
			want:    user.User{},
			wantErr: exceptions.NewInvalidArgumentError("Password", "password must not empty", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ValidateUsernamePassword{
				userRepo: tt.dependencies.userRepoFunc(t),
			}
			got, err := h.Handle(context.TODO(), tt.args)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
