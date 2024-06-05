package propre

import (
	"context"
	"net/http"
)

type HTTPSendable interface {
	ContentType() string
	Encode() ([]byte, error)
	StatusCode() int
}

type HTTPResponse[View HTTPSendable, Writer http.ResponseWriter] struct {
	headers              http.Header
	genericInternalError []byte
}

type HTTPResponseOpts[View HTTPSendable, Writer http.ResponseWriter] func(r *HTTPResponse[View, Writer])

func WithHTTPResponseHeaders[View HTTPSendable, Writer http.ResponseWriter](headers http.Header) HTTPResponseOpts[View, Writer] {
	return func(r *HTTPResponse[View, Writer]) {
		r.headers = headers
	}
}

func NewHTTPResponse[View HTTPSendable, Writer http.ResponseWriter](opts ...HTTPResponseOpts[View, Writer]) *HTTPResponse[View, Writer] {
	response := &HTTPResponse[View, Writer]{}
	for _, opt := range opts {
		opt(response)
	}

	return response
}

func (r *HTTPResponse[View, Writer]) Send(ctx context.Context, rw http.ResponseWriter, data View) {
	rw.Header().Set("content-type", data.ContentType())

	if len(r.headers) > 0 {
		for header, values := range r.headers {
			for _, headerValue := range values {
				rw.Header().Set(header, headerValue)
			}
		}
	}

	encoded, err := data.Encode()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(r.genericInternalError)
		return
	}

	rw.WriteHeader(data.StatusCode())
	rw.Write(encoded)
}
