package validate_user_token

import (
	"testing"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserToken_Handle(t *testing.T) {
	tokenCreatedAt := time.Now().Add(-10 * time.Second).UnixMilli()
	tokenExpiredAtInPast := time.Now().Add(-1 * time.Second).UnixMilli()
	tokenExpiredAtInFeature := time.Now().Add(10 * time.Minute).UnixMilli()
	type dependencies struct {
		validatorFunc func(tt *testing.T) tokenValidator
	}
	type args struct {
		token string
	}
	tests := []struct {
		name         string
		dependencies dependencies
		args         args
		want         authenticate.UserAuthenticateData
		wantErr      error
	}{
		{
			name: "empty token",
			dependencies: dependencies{
				validatorFunc: func(tt *testing.T) tokenValidator {
					m := NewMocktokenValidator(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				token: "",
			},
			want:    authenticate.UserAuthenticateData{},
			wantErr: exceptions.NewPreconditionError(exceptions.PreconditionTypeInvalidToken, exceptions.SubjectAuthentication, "token must not empty", nil),
		},
		{
			name: "token valid",
			dependencies: dependencies{
				validatorFunc: func(tt *testing.T) tokenValidator {
					m := NewMocktokenValidator(gomock.NewController(tt))
					m.EXPECT().ValidateToken("valid_token").Return(authenticate.UserAuthenticateData{
						UserID:       1,
						Username:     "banhdanh",
						UserFullName: "Banh Quoc Danh",
						CreatedAt:    tokenCreatedAt,
						ExpireAt:     tokenExpiredAtInFeature,
					}, nil)
					return m
				},
			},
			args: args{
				token: "valid_token",
			},
			want: authenticate.UserAuthenticateData{
				UserID:       1,
				Username:     "banhdanh",
				UserFullName: "Banh Quoc Danh",
				CreatedAt:    tokenCreatedAt,
				ExpireAt:     tokenExpiredAtInFeature,
			},
			wantErr: nil,
		},
		{
			name: "token expired",
			dependencies: dependencies{
				validatorFunc: func(tt *testing.T) tokenValidator {
					m := NewMocktokenValidator(gomock.NewController(tt))
					m.EXPECT().ValidateToken("valid_token").Return(authenticate.UserAuthenticateData{
						UserID:       1,
						Username:     "banhdanh",
						UserFullName: "Banh Quoc Danh",
						CreatedAt:    tokenCreatedAt,
						ExpireAt:     tokenExpiredAtInPast,
					}, nil)
					return m
				},
			},
			args: args{
				token: "valid_token",
			},
			want: authenticate.UserAuthenticateData{},
			wantErr: exceptions.NewPreconditionError(
				exceptions.PreconditionTypeTokenExpired,
				exceptions.SubjectAuthentication,
				"token expired",
				map[string]interface{}{
					"expired_at": tokenExpiredAtInPast,
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := ValidateUserToken{
				validator: tt.dependencies.validatorFunc(t),
			}
			got, err := v.Handle(nil, tt.args.token)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
