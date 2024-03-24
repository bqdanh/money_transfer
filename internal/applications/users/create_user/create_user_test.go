package create_user

import (
	"context"
	"fmt"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
	matcher_gomock "github.com/bqdanh/money_transfer/pkg/gomock_matcher"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateUser(t *testing.T) {
	_, err := NewCreateUser(nil)
	assert.EqualError(t, err, "user repository must not nil")

	m := NewMockuserRepository(gomock.NewController(t))
	h, err := NewCreateUser(m)
	assert.NoError(t, err)
	assert.Equal(t, h.ur, m)
}

func TestCreateUser_Handle(t *testing.T) {
	mockRepoError := fmt.Errorf("mock repository error")
	type dependencies struct {
		urf func(tt *testing.T) userRepository
	}
	type expected struct {
		u   user.User
		err error
	}

	testcases := []struct {
		name         string
		dependencies dependencies
		params       CreateUserParams
		expected     expected
	}{
		{
			name: "happy case",
			dependencies: dependencies{
				urf: func(tt *testing.T) userRepository {
					ur := NewMockuserRepository(gomock.NewController(tt))
					ur.EXPECT().CreateUser(gomock.AssignableToTypeOf(context.TODO()), matcher_gomock.MatchesFunc(func(got interface{}) error {
						v, ok := got.(user.User)
						if !ok {
							return fmt.Errorf("want(%T), but got(%T)", user.User{}, got)
						}
						if v.ID != 0 {
							return fmt.Errorf("want 0, but got (%v)", v.ID)
						}
						if v.Password == "pw_test" {
							return fmt.Errorf("must hash password")
						}
						return nil
					})).Return(user.User{
						ID:       1,
						UserName: "username_test",
						Password: "$2a$10$NCFMu48aCOFEX9TU2Nut0O9TrmR6TRJREBALxUHsx9Ht2cq3PTu7",
						FullName: "xyz",
						Phone:    "8400000",
					}, nil)
					return ur
				},
			},
			params: CreateUserParams{
				UserName: "username_test",
				Password: "pw_test",
				FullName: "name test",
				Phone:    "8400000",
			},
			expected: expected{
				u: user.User{
					ID:       1,
					UserName: "username_test",
					Password: "$2a$10$NCFMu48aCOFEX9TU2Nut0O9TrmR6TRJREBALxUHsx9Ht2cq3PTu7",
					FullName: "xyz",
					Phone:    "8400000",
				},
				err: nil,
			},
		},
		{
			name: "user repository got error",
			dependencies: dependencies{
				urf: func(tt *testing.T) userRepository {
					ur := NewMockuserRepository(gomock.NewController(tt))
					ur.EXPECT().CreateUser(gomock.AssignableToTypeOf(context.TODO()), matcher_gomock.MatchesFunc(func(got interface{}) error {
						v, ok := got.(user.User)
						if !ok {
							return fmt.Errorf("want(%T), but got(%T)", user.User{}, got)
						}
						if v.ID != 0 {
							return fmt.Errorf("want 0, but got (%v)", v.ID)
						}
						if v.Password == "pw_test" {
							return fmt.Errorf("must hash password")
						}
						return nil
					})).Return(user.User{}, mockRepoError)
					return ur
				},
			},
			params: CreateUserParams{
				UserName: "username_test",
				Password: "pw_test",
				FullName: "name test",
				Phone:    "8400000",
			},
			expected: expected{
				u:   user.User{},
				err: mockRepoError,
			},
		},
		{
			name: "invalid user name",
			dependencies: dependencies{
				urf: func(tt *testing.T) userRepository {
					ur := NewMockuserRepository(gomock.NewController(tt))
					return ur
				},
			},
			params: CreateUserParams{
				UserName: "",
				Password: "pw_test",
				FullName: "name test",
				Phone:    "8400000",
			},
			expected: expected{
				u:   user.User{},
				err: exceptions.NewInvalidArgumentError("UserName", "length of user name must between 6 and 16", map[string]interface{}{}),
			},
		},
		{
			name: "invalid password",
			dependencies: dependencies{
				urf: func(tt *testing.T) userRepository {
					ur := NewMockuserRepository(gomock.NewController(tt))
					return ur
				},
			},
			params: CreateUserParams{
				UserName: "username_test",
				Password: "",
				FullName: "name test",
				Phone:    "8400000",
			},
			expected: expected{
				u:   user.User{},
				err: exceptions.NewInvalidArgumentError("Password", "length of password must between 6 and 16", map[string]interface{}{}),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t2 *testing.T) {
			h, err := NewCreateUser(tc.dependencies.urf(t2))
			assert.NoError(t2, err)
			u, err := h.Handle(context.TODO(), tc.params)
			assert.ErrorIs(t2, err, tc.expected.err)
			assert.EqualValues(t2, tc.expected.u, u)
		})
	}
}
