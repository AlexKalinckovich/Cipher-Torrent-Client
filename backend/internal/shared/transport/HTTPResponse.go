package transport

type HTTPResponse struct {
	status int
	Body   any `json:"message"`
}

func NewHTTPResponse(status int, message any) HTTPResponse {
	return HTTPResponse{status: status, Body: message}
}

func (r HTTPResponse) StatusCode() int {
	return r.status
}
