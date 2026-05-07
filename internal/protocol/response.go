package protocol

type Response[T any] struct {
	Data  *T             `json:"data"`
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}
