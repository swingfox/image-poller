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

type ImageResponse struct {
	Limit     int         `json:"limit"`
	ImageData []ImageData `json:"data"`
}

type ImageData struct {
	ID   string `json:"id"`
	Hits int    `json:"hits"`
	Uri  string `json:"uri"`
}

type Photo struct {
	Src PhotosSrc `json:"src"`
}

type PhotosSrc struct {
	Original string `json:"original"`
}

type ProviderResponse struct {
	Page    int     `json:"page"`
	PerPage int     `json:"per_page"`
	Photos  []Photo `json:"photos"`
}
