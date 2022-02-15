package error

import "errors"

var (
	ServerErr            = errors.New("server error")
	EmptyBodyRequest     = errors.New("empty body request")
	InvalidUrl           = errors.New("invalid url")
	ErrorNoQueryParam    = errors.New("query param does not exists")
	ErrorExpiredLink     = errors.New("expired link")
)
