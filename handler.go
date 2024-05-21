package propre

import "net/http"

type HTTPHandler[Input, Output any] struct {
	requestDecoder RequestDecoder[Input]
	useCaseHandler UseCaseHandler[Input, Output]
	responseSender HTTPResponseSender[Output]
}

func NewHTTPHandler[Input, Output any](
	requestDecoder RequestDecoder[Input],
	useCaseHandler UseCaseHandler[Input, Output],
	responseSender HTTPResponseSender[Output],
) *HTTPHandler[Input, Output] {
	return &HTTPHandler[Input, Output]{
		requestDecoder: requestDecoder,
		useCaseHandler: useCaseHandler,
		responseSender: responseSender,
	}
}

func (handler *HTTPHandler[Input, Output]) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	input := handler.requestDecoder.Decode(req)
	output := handler.useCaseHandler.Handle(req.Context(), input)
	handler.responseSender.Send(req.Context(), rw, output)
}
