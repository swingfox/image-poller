package user

type Request struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}
