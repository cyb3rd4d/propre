package propre

import (
	"context"
	"net/http"
)

var (
	defaultInternalError = []byte("internal error")
)

type HTTPSendable interface {
	ContentType() string
	Encode() ([]byte, error)
	StatusCode() int
}

type HTTPResponse[View HTTPSendable] struct {
	headers              http.Header
	genericInternalError []byte
}

type HTTPResponseOpts[View HTTPSendable] func(r *HTTPResponse[View])

func WithHTTPResponseHeaders[View HTTPSendable](headers http.Header) HTTPResponseOpts[View] {
	return func(r *HTTPResponse[View]) {
		r.headers = headers
	}
}

func WithGenericInternalError[View HTTPSendable](payload []byte) HTTPResponseOpts[View] {
	return func(r *HTTPResponse[View]) {
		r.genericInternalError = payload
	}
}

func NewHTTPResponse[View HTTPSendable](opts ...HTTPResponseOpts[View]) *HTTPResponse[View] {
	response := &HTTPResponse[View]{}
	for _, opt := range opts {
		opt(response)
	}

	return response
}

func (r *HTTPResponse[View]) Send(ctx context.Context, rw http.ResponseWriter, data View) {
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
		internalError := defaultInternalError
		if r.genericInternalError != nil {
			internalError = r.genericInternalError
		}

		rw.Write(internalError)
		return
	}

	rw.WriteHeader(data.StatusCode())
	rw.Write(encoded)
}
