package net

import (
	"net/http"
)

//CustomResponse is a helper for handles the custom response object
type CustomResponse struct {
	Body       []byte
	Headers    http.Header
	StatusCode int
}

// Write writes the data to the connection as part of an HTTP reply.
func (c *CustomResponse) Write(b []byte) (int, error) {
	c.Body = b

	if c.StatusCode == 0 {
		c.StatusCode = 200
	}
	return 0, nil
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (c *CustomResponse) WriteHeader(s int) {
	if s != 0 {
		c.StatusCode = s
	}
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (c *CustomResponse) Header() http.Header {
	c.Headers = http.Header{}
	return c.Headers
}
