package response_types

type ApiResponse struct {
	Content    any   `json:"content,omitempty"`
	Error      error `json:"error,omitempty"`
	StatusCode int   `json:"-"`
}
