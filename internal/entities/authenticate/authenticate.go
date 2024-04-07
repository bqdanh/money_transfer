package authenticate

type UserAuthenticateData struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	UserFullName string `json:"user_full_name"`
	CreatedAt    int64  `json:"created_at"` //unix millisecond timestamp
	ExpireAt     int64  `json:"expire_at"`  //unix millisecond timestamp
}
