package transport

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
)

type HTTPResponse struct {
	status int
	Body   common.ApiError
}

func NewHTTPResponse(status int, body common.ApiError) HTTPResponse {
	return HTTPResponse{status: status, Body: body}
}

func (r HTTPResponse) StatusCode() int {
	return r.status
}
