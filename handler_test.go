package propre_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyb3rd4d/propre"
	"github.com/stretchr/testify/mock"
)

func TestHTTPHandlerImplementsHTTPHandler(t *testing.T) {
	handler := propre.NewHTTPHandler(
		new(requestDecoderMock[any]),
		new(useCaseHandlerMock[any, any]),
		new(httpResponseSenderMock[any]),
	)

	f := func(h http.Handler) {}
	f(handler)
}

func TestHTTPHandlerUsesARequestDecoderThenAUseCaseHandlerThenSendsTheResponse(t *testing.T) {
	requestDecoder := new(requestDecoderMock[any])
	defer requestDecoder.AssertExpectations(t)

	useCaseHandler := new(useCaseHandlerMock[any, any])
	defer useCaseHandler.AssertExpectations(t)

	responseSender := new(httpResponseSenderMock[any])
	defer responseSender.AssertExpectations(t)

	handler := propre.NewHTTPHandler(requestDecoder, useCaseHandler, responseSender)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	useCaseInput := "some input"
	useCaseOutput := "some output"

	requestDecoder.On("Decode", req).Return(useCaseInput, nil)

	ctxArgMatcher := mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})

	useCaseHandler.On("Handle", ctxArgMatcher, useCaseInput).Return(useCaseOutput)
	responseSender.On("Send", ctxArgMatcher, rw, useCaseOutput)

	handler.ServeHTTP(rw, req)
}

type requestDecoderMock[Input any] struct {
	mock.Mock
}

func (m *requestDecoderMock[Input]) Decode(req *http.Request) Input {
	args := m.Called(req)

	return args.Get(0).(Input)
}

type useCaseHandlerMock[Input, Output any] struct {
	mock.Mock
}

func (m *useCaseHandlerMock[Input, Output]) Handle(ctx context.Context, input Input) Output {
	return m.Called(ctx, input).Get(0).(Output)
}

type httpResponseSenderMock[Output any] struct {
	mock.Mock
}

func (m *httpResponseSenderMock[Output]) Send(ctx context.Context, rw http.ResponseWriter, output Output) {
	m.Called(ctx, rw, output)
}
