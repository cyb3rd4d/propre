package propre

import (
	"context"
	"net/http"
)

var (
	defaultInternalError = []byte("internal error")
)

// HTTPSendable is the constraint that view models must implement to be
// used in [HTTPResponse].
type HTTPSendable interface {
	ContentType(context.Context) string
	Encode(context.Context) ([]byte, error)
	StatusCode(context.Context) int
}

// HTTPResponse sends the final response to the client.
// It needs a View type parameter implementing [HTTPSendable] to know
// the content type and the status code to return, and it delegates
// the payload encoding to the view model.
type HTTPResponse[View HTTPSendable] struct {
	headers              http.Header
	genericInternalError []byte
}

// HTTPResponseOpts is the alias for the [HTTPResponse] builder options.
type HTTPResponseOpts[View HTTPSendable] func(r *HTTPResponse[View])

// WithHTTPResponseHeaders is an [HTTPResponse] option to set common headers
// for every single response.
func WithHTTPResponseHeaders[View HTTPSendable](headers http.Header) HTTPResponseOpts[View] {
	return func(r *HTTPResponse[View]) {
		r.headers = headers
	}
}

// WithGenericInternalError is an [HTTPResponse] option to define a custom
// payload for internal errors.
// Such errors can be raised if the view model encoding fails.
func WithGenericInternalError[View HTTPSendable](payload []byte) HTTPResponseOpts[View] {
	return func(r *HTTPResponse[View]) {
		r.genericInternalError = payload
	}
}

// NewHTTPResponse returns an [HTTPResponse]. [HTTPResponseOpts] can be passed
// to customize the common response headers and the default internal error payload.
func NewHTTPResponse[View HTTPSendable](opts ...HTTPResponseOpts[View]) *HTTPResponse[View] {
	response := &HTTPResponse[View]{}
	for _, opt := range opts {
		opt(response)
	}

	return response
}

// Send returns the response to the client.
// It first sets the content type and the common headers if some have been defined
// and the status code, and the payload is encoded and is sent through the [http.ResponseWriter].
func (r *HTTPResponse[View]) Send(ctx context.Context, rw http.ResponseWriter, data View) {
	rw.Header().Set("content-type", data.ContentType(ctx))

	if len(r.headers) > 0 {
		for header, values := range r.headers {
			for _, headerValue := range values {
				rw.Header().Set(header, headerValue)
			}
		}
	}

	encoded, err := data.Encode(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		internalError := defaultInternalError
		if r.genericInternalError != nil {
			internalError = r.genericInternalError
		}

		rw.Write(internalError)
		return
	}

	rw.WriteHeader(data.StatusCode(ctx))
	rw.Write(encoded)
}
