package response_types

type Responder interface {
	GetStatusCode() int
	GetError() error
}

type ApiResponse struct {
	Content    any   `json:"content,omitempty"`
	Error      error `json:"error,omitempty"`
	StatusCode int   `json:"-"`
}

type FileResponse struct {
	Name       string
	StatusCode int
	Error      error
}

func (r ApiResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r ApiResponse) GetError() error {
	return r.Error
}

func (r FileResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r FileResponse) GetError() error {
	return r.Error
}
