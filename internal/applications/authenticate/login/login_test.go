package login

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bqdanh/money_transfer/internal/applications/authenticate/generate_user_token"
	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
	"github.com/bqdanh/money_transfer/pkg/gomock_matcher"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Handle(t *testing.T) {
	mockUser := user.User{
		ID:       1,
		UserName: "banhdanh",
		Password: "$2a$10$lepW1ozns2uifq0SNs.sduKceRysA3WvhL/puF.mMlDXMHdMprKza",
		FullName: "Banh Quoc Danh",
		Phone:    "0973280276",
	}
	generateErr := fmt.Errorf("generate token error")
	type fields struct {
		userNamePasswordValidatorFunc func(tt *testing.T) userNamePasswordValidator
		generateTokenConfig           generate_user_token.Config
		mockTokenGeneratorFunc        func(tt *testing.T) *generate_user_token.MocktokenGenerator
	}

	tests := []struct {
		name    string
		fields  fields
		args    LoginParams
		want    LoginResponse
		wantErr error
	}{
		{
			name: "login success",
			fields: fields{
				userNamePasswordValidatorFunc: func(tt *testing.T) userNamePasswordValidator {
					m := NewMockuserNamePasswordValidator(gomock.NewController(tt))
					m.EXPECT().ValidateUserNamePassword(gomock.Any(), "banhdanh", "123").Return(mockUser, nil)
					return m
				},
				generateTokenConfig: generate_user_token.Config{
					TokenDuration: 1 * time.Hour,
				},
				mockTokenGeneratorFunc: func(tt *testing.T) *generate_user_token.MocktokenGenerator {
					m := generate_user_token.NewMocktokenGenerator(gomock.NewController(tt))
					m.EXPECT().GenerateToken(gomock_matcher.MatchesFunc(func(got interface{}) error {
						v, ok := got.(authenticate.UserAuthenticateData)
						if !ok {
							return fmt.Errorf("expect authenticate.UserAuthenticateData, got %T", got)
						}
						if v.UserID != mockUser.ID {
							return fmt.Errorf("expect UserID %d, got %d", mockUser.ID, v.UserID)
						}
						if v.Username != mockUser.UserName {
							return fmt.Errorf("expect Username %s, got %s", mockUser.UserName, v.Username)
						}
						if v.UserFullName != mockUser.FullName {
							return fmt.Errorf("expect UserFullName %s, got %s", mockUser.FullName, v.UserFullName)
						}
						if time.Now().UnixMilli()-v.CreatedAt > int64(10*time.Second/time.Millisecond) {
							return fmt.Errorf("expect CreatedAt %d~10s, got %d", time.Now().UnixMilli(), v.CreatedAt)
						}
						if v.ExpireAt < time.Now().UnixMilli() {
							return fmt.Errorf("expect ExpireAt > %d, got %d", time.Now().UnixMilli(), v.ExpireAt)
						}
						if time.Now().Add(1*time.Hour).UnixMilli()-v.ExpireAt > int64(10*time.Second/time.Millisecond) {
							return fmt.Errorf("expect ExpireAt ~ %d, got %d", time.Now().Add(1*time.Hour).UnixMilli(), v.ExpireAt)
						}
						return nil
					})).Return("token_xyz", nil)
					return m
				},
			},
			args: LoginParams{
				Username: "banhdanh",
				Password: "123",
			},
			want: LoginResponse{
				Token: "token_xyz",
			},
			wantErr: nil,
		},
		{
			name: "login with empty password",
			fields: fields{
				userNamePasswordValidatorFunc: func(tt *testing.T) userNamePasswordValidator {
					m := NewMockuserNamePasswordValidator(gomock.NewController(tt))
					return m
				},
				generateTokenConfig: generate_user_token.Config{
					TokenDuration: 1 * time.Hour,
				},
				mockTokenGeneratorFunc: func(tt *testing.T) *generate_user_token.MocktokenGenerator {
					m := generate_user_token.NewMocktokenGenerator(gomock.NewController(tt))
					return m
				},
			},
			args: LoginParams{
				Username: "banhdanh",
				Password: "",
			},
			want: LoginResponse{
				Token: "",
			},
			wantErr: exceptions.NewInvalidArgumentError("Password", "password must not empty", nil),
		},
		{
			name: "login with empty username",
			fields: fields{
				userNamePasswordValidatorFunc: func(tt *testing.T) userNamePasswordValidator {
					m := NewMockuserNamePasswordValidator(gomock.NewController(tt))
					return m
				},
				generateTokenConfig: generate_user_token.Config{
					TokenDuration: 1 * time.Hour,
				},
				mockTokenGeneratorFunc: func(tt *testing.T) *generate_user_token.MocktokenGenerator {
					m := generate_user_token.NewMocktokenGenerator(gomock.NewController(tt))
					return m
				},
			},
			args: LoginParams{
				Username: "",
				Password: "123",
			},
			want: LoginResponse{
				Token: "",
			},
			wantErr: exceptions.NewInvalidArgumentError("Username", "username must not empty", nil),
		},
		{
			name: "login with fail validate",
			fields: fields{
				userNamePasswordValidatorFunc: func(tt *testing.T) userNamePasswordValidator {
					m := NewMockuserNamePasswordValidator(gomock.NewController(tt))
					m.EXPECT().ValidateUserNamePassword(gomock.Any(), "banhdanh", "123").
						Return(
							user.User{},
							exceptions.NewPreconditionError(
								exceptions.PreconditionTypePasswordNotMatch,
								exceptions.SubjectUser,
								"password not match",
								map[string]interface{}{},
							),
						)
					return m
				},
				generateTokenConfig: generate_user_token.Config{
					TokenDuration: 1 * time.Hour,
				},
				mockTokenGeneratorFunc: func(tt *testing.T) *generate_user_token.MocktokenGenerator {
					m := generate_user_token.NewMocktokenGenerator(gomock.NewController(tt))
					return m
				},
			},
			args: LoginParams{
				Username: "banhdanh",
				Password: "123",
			},
			want: LoginResponse{
				Token: "",
			},
			wantErr: exceptions.NewPreconditionError(
				exceptions.PreconditionTypePasswordNotMatch,
				exceptions.SubjectUser,
				"password not match",
				map[string]interface{}{},
			),
		},
		{
			name: "login with fail generate token",
			fields: fields{
				userNamePasswordValidatorFunc: func(tt *testing.T) userNamePasswordValidator {
					m := NewMockuserNamePasswordValidator(gomock.NewController(tt))
					m.EXPECT().ValidateUserNamePassword(gomock.Any(), "banhdanh", "123").Return(mockUser, nil)
					return m
				},
				generateTokenConfig: generate_user_token.Config{
					TokenDuration: 1 * time.Hour,
				},
				mockTokenGeneratorFunc: func(tt *testing.T) *generate_user_token.MocktokenGenerator {
					m := generate_user_token.NewMocktokenGenerator(gomock.NewController(tt))
					m.EXPECT().GenerateToken(gomock_matcher.MatchesFunc(func(got interface{}) error {
						v, ok := got.(authenticate.UserAuthenticateData)
						if !ok {
							return fmt.Errorf("expect authenticate.UserAuthenticateData, got %T", got)
						}
						if v.UserID != mockUser.ID {
							return fmt.Errorf("expect UserID %d, got %d", mockUser.ID, v.UserID)
						}
						if v.Username != mockUser.UserName {
							return fmt.Errorf("expect Username %s, got %s", mockUser.UserName, v.Username)
						}
						if v.UserFullName != mockUser.FullName {
							return fmt.Errorf("expect UserFullName %s, got %s", mockUser.FullName, v.UserFullName)
						}
						if time.Now().UnixMilli()-v.CreatedAt > int64(10*time.Second/time.Millisecond) {
							return fmt.Errorf("expect CreatedAt %d~10s, got %d", time.Now().UnixMilli(), v.CreatedAt)
						}
						if v.ExpireAt < time.Now().UnixMilli() {
							return fmt.Errorf("expect ExpireAt > %d, got %d", time.Now().UnixMilli(), v.ExpireAt)
						}
						if time.Now().Add(1*time.Hour).UnixMilli()-v.ExpireAt > int64(10*time.Second/time.Millisecond) {
							return fmt.Errorf("expect ExpireAt ~ %d, got %d", time.Now().Add(1*time.Hour).UnixMilli(), v.ExpireAt)
						}
						return nil
					})).Return("", generateErr)
					return m
				},
			},
			args: LoginParams{
				Username: "banhdanh",
				Password: "123",
			},
			want: LoginResponse{
				Token: "",
			},
			wantErr: generateErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator, err := generate_user_token.NewGenerateUserToken(tt.fields.mockTokenGeneratorFunc(t), tt.fields.generateTokenConfig)
			assert.NoError(t, err)

			l := Login{
				userNamePasswordValidator: tt.fields.userNamePasswordValidatorFunc(t),
				generateUserToken:         generator,
			}
			got, err := l.Handle(context.TODO(), tt.args)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
