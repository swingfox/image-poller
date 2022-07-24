package image

type Request struct {
	URI   string `validate:"required"`
	Owner string
}
type Response struct {
	ID   string
	Hits int
	URI  string
}
