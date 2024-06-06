package propre

import "net/http"

// HTTPHandler is the main component to handle HTTP requests with propre.
// Each endpoint requires:
//   - a request decoder to transform an HTTP request to a use case input,
//   - a use case handler responsible of checking possible input errors, storing data in a repository,
//     writing business code,
//   - a presenter to send final HTTP response depending on the output returned by the use case.
//
// It implements http.Handler to be easily used by any HTTP "ServeMux".
type HTTPHandler[Input, Output any] struct {
	requestDecoder RequestDecoder[Input]
	useCaseHandler UseCaseHandler[Input, Output]
	presenter      Presenter[Output, http.ResponseWriter]
}

// NewHTTPHandler builds an HTTPHandler with the given dependencies.
func NewHTTPHandler[Input, Output any](
	requestDecoder RequestDecoder[Input],
	useCaseHandler UseCaseHandler[Input, Output],
	presenter Presenter[Output, http.ResponseWriter],
) *HTTPHandler[Input, Output] {
	return &HTTPHandler[Input, Output]{
		requestDecoder: requestDecoder,
		useCaseHandler: useCaseHandler,
		presenter:      presenter,
	}
}

// ServeHTTP allows HTTPHandler to be used by any HTTP "ServeMux"
// Its implementation is straightforward:
//   - a request decoder transforms an HTTP request to a use case input,
//   - a use case handler takes the previous input to handle the business logic,
//   - a presenter sends the final HTTP response depending on the output returned by the use case.
func (handler *HTTPHandler[Input, Output]) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	input := handler.requestDecoder.Decode(req)
	output := handler.useCaseHandler.Handle(req.Context(), input)
	handler.presenter.Present(req.Context(), rw, output)
}
