package user

type User struct {
	ID int64
	//TODO: move UserName and Password to Auth entity
	UserName string
	Password string

	Email    string
	FullName string
	Phone    string
}
