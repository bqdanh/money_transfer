package generate_user_token

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
	"github.com/bqdanh/money_transfer/pkg/gomock_matcher"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerateUserToken(t *testing.T) {
	type args struct {
		g tokenGenerator
		d time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    GenerateUserToken
		wantErr error
	}{
		{
			name: "nil generator",
			args: args{
				g: nil,
				d: 15 * time.Minute,
			},
			want:    GenerateUserToken{},
			wantErr: exceptions.NewInvalidArgumentError("generator", "generator must not nil", nil),
		},
		{
			name: "invalid duration",
			args: args{
				g: NewMocktokenGenerator(gomock.NewController(t)),
				d: 0,
			},
			want:    GenerateUserToken{},
			wantErr: exceptions.NewInvalidArgumentError("duration", "duration must greater than 0", nil),
		},
		{
			name: "valid",
			args: args{
				g: NewMocktokenGenerator(gomock.NewController(t)),
				d: 15 * time.Minute,
			},
			want: GenerateUserToken{
				generator:     NewMocktokenGenerator(gomock.NewController(t)),
				tokenDuration: 15 * time.Minute,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenerateUserToken(tt.args.g, Config{TokenDuration: tt.args.d})
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateUserToken_Handle(t *testing.T) {
	tokenDuration := 15 * time.Minute
	errTest := fmt.Errorf("test error")

	type dependencies struct {
		tokenDuration time.Duration
		generatorFunc func(tt *testing.T) tokenGenerator
	}
	type args struct {
		u user.User
	}
	tests := []struct {
		name         string
		dependencies dependencies
		args         args
		want         string
		wantErr      error
	}{
		{
			name: "generate user token",
			dependencies: dependencies{
				tokenDuration: tokenDuration,
				generatorFunc: func(tt *testing.T) tokenGenerator {
					m := NewMocktokenGenerator(gomock.NewController(tt))
					m.EXPECT().GenerateToken(gomock_matcher.MatchesFunc(func(got interface{}) error {
						v, ok := got.(authenticate.UserAuthenticateData)
						if !ok {
							return fmt.Errorf("want(%T), but got(%T)", authenticate.UserAuthenticateData{}, got)
						}
						if v.UserID != 1 {
							return fmt.Errorf("want (1), but got (%v)", v.UserID)
						}
						if v.Username != "banhdanh" {
							return fmt.Errorf("want banhdanh, but got (%v)", v.Username)
						}
						if v.UserFullName != "Banh Quoc Danh" {
							return fmt.Errorf("want (Banh Quoc Danh), but got (%v)", v.UserFullName)
						}
						//validate token created at near current time, expected different less than 10s
						if time.Now().UnixMilli()-v.CreatedAt > int64(10*time.Second/time.Millisecond) {
							return fmt.Errorf("want (~%d), but got (%v)", time.Now().UnixMilli(), v.CreatedAt)
						}
						if time.Now().UnixMilli() < v.CreatedAt {
							return fmt.Errorf("want < (%d), but got (%v)", time.Now().UnixMilli(), v.CreatedAt)
						}
						if v.ExpireAt <= time.Now().UnixMilli() {
							return fmt.Errorf("want > (%d), but got (%v)", time.Now().UnixMilli(), v.ExpireAt)
						}
						if v.ExpireAt > time.Now().Add(tokenDuration).UnixMilli() {
							return fmt.Errorf("want < (%d), but got (%v)", time.Now().Add(tokenDuration).UnixMilli(), v.ExpireAt)
						}
						return nil
					})).Return("token_test", nil)
					return m
				},
			},
			args: args{
				u: user.User{
					ID:       1,
					UserName: "banhdanh",
					Password: "",
					FullName: "Banh Quoc Danh",
					Phone:    "0973280276",
				},
			},
			want:    "token_test",
			wantErr: nil,
		},
		{
			name: "generate user token fail",
			dependencies: dependencies{
				tokenDuration: tokenDuration,
				generatorFunc: func(tt *testing.T) tokenGenerator {
					m := NewMocktokenGenerator(gomock.NewController(tt))
					m.EXPECT().GenerateToken(gomock_matcher.MatchesFunc(func(got interface{}) error {
						v, ok := got.(authenticate.UserAuthenticateData)
						if !ok {
							return fmt.Errorf("want(%T), but got(%T)", authenticate.UserAuthenticateData{}, got)
						}
						if v.UserID != 1 {
							return fmt.Errorf("want (1), but got (%v)", v.UserID)
						}
						if v.Username != "banhdanh" {
							return fmt.Errorf("want banhdanh, but got (%v)", v.Username)
						}
						if v.UserFullName != "Banh Quoc Danh" {
							return fmt.Errorf("want (Banh Quoc Danh), but got (%v)", v.UserFullName)
						}
						//validate token created at near current time, expected different less than 10s
						if time.Now().UnixMilli()-v.CreatedAt > int64(10*time.Second/time.Millisecond) {
							return fmt.Errorf("want (~%d), but got (%v)", time.Now().UnixMilli(), v.CreatedAt)
						}
						if time.Now().UnixMilli() < v.CreatedAt {
							return fmt.Errorf("want < (%d), but got (%v)", time.Now().UnixMilli(), v.CreatedAt)
						}
						if v.ExpireAt <= time.Now().UnixMilli() {
							return fmt.Errorf("want > (%d), but got (%v)", time.Now().UnixMilli(), v.ExpireAt)
						}
						if v.ExpireAt > time.Now().Add(tokenDuration).UnixMilli() {
							return fmt.Errorf("want < (%d), but got (%v)", time.Now().Add(tokenDuration).UnixMilli(), v.ExpireAt)
						}
						return nil
					})).Return("", errTest)
					return m
				},
			},
			args: args{
				u: user.User{
					ID:       1,
					UserName: "banhdanh",
					Password: "",
					FullName: "Banh Quoc Danh",
					Phone:    "0973280276",
				},
			},
			want:    "",
			wantErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GenerateUserToken{
				generator:     tt.dependencies.generatorFunc(t),
				tokenDuration: tt.dependencies.tokenDuration,
			}
			got, err := g.Handle(context.TODO(), tt.args.u)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
