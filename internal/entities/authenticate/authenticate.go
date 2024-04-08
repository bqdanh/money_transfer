package authenticate

import "context"

type UserAuthenticateData struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	UserFullName string `json:"user_full_name"`
	CreatedAt    int64  `json:"created_at"` //unix millisecond timestamp
	ExpireAt     int64  `json:"expire_at"`  //unix millisecond timestamp
}

type userAuthenticateDataKey struct{}

func ContextWithAuthentication(ctx context.Context, data UserAuthenticateData) context.Context {
	return context.WithValue(ctx, userAuthenticateDataKey{}, data)
}

func AuthenticationFromContext(ctx context.Context) UserAuthenticateData {
	vl := ctx.Value(userAuthenticateDataKey{})
	if vl == nil {
		return UserAuthenticateData{}
	}
	data, ok := vl.(UserAuthenticateData)
	if !ok {
		return UserAuthenticateData{}
	}
	return data
}
