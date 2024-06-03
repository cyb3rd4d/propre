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
	headers              []string
	genericInternalError []byte
}

// TODO: add tests
type HTTPResponseOpts func(r *HTTPResponse[HTTPSendable, http.ResponseWriter])

func WithHTTPResponseHeaders(headers []string) HTTPResponseOpts {
	return func(r *HTTPResponse[HTTPSendable, http.ResponseWriter]) {
		r.headers = headers
	}
}

func NewHTTPResponse[View HTTPSendable, Writer http.ResponseWriter]() *HTTPResponse[View, Writer] {
	return &HTTPResponse[View, Writer]{}
}

func (r *HTTPResponse[View, Writer]) Send(ctx context.Context, rw http.ResponseWriter, data View) {
	rw.Header().Set("content-type", data.ContentType())
	encoded, err := data.Encode()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(r.genericInternalError)
		return
	}

	rw.WriteHeader(data.StatusCode())
	rw.Write(encoded)
}
