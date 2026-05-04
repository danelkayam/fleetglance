package protocol

type Response struct {
	Data  any            `json:"data"`
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}
