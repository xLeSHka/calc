package application

import "errors"

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrCanNotEncodeResp = errors.New("can not encode resp")
)
