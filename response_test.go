package propre_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyb3rd4d/propre"
)

type payload[T any] struct {
	OK    *T              `json:"ok,omitempty"`
	Error *errorViewModel `json:"error,omitempty"`
}

func (p payload[T]) ContentType() string {
	return "application/json"
}

func (p payload[T]) Encode() ([]byte, error) {
	return json.Marshal(p)
}

func (p payload[T]) StatusCode() int {
	s := http.StatusOK
	if p.Error != nil {
		s = http.StatusInternalServerError
	}

	return s
}

type errorViewModel struct {
	Message string `json:"message"`
}

type okViewModel struct {
	Data string `json:"data"`
}

type outputData struct {
	SomeField string
}

type monad struct {
	Data  outputData
	Error error
}

type testPresenter[Output monad, Writer http.ResponseWriter] struct {
	response *propre.HTTPResponse[payload[okViewModel], Writer]
}

func (s *testPresenter[Output, Writer]) Present(ctx context.Context, rw http.ResponseWriter, output monad) {
	var p payload[okViewModel]

	if output.Error != nil {
		p.Error = &errorViewModel{
			Message: fmt.Sprintf("an error occurred: %s", output.Error),
		}

		s.response.Send(ctx, rw, p)
		return
	}

	p.OK = &okViewModel{
		Data: fmt.Sprintf("success: %s", output.Data.SomeField),
	}

	s.response.Send(ctx, rw, p)
}

type responseTestCase struct {
	output               monad
	expectedHTTPStatus   int
	expectedJSONResponse []byte
}

func TestResponse(t *testing.T) {
	presenter := &testPresenter[monad, http.ResponseWriter]{
		response: propre.NewHTTPResponse[payload[okViewModel], http.ResponseWriter](),
	}

	testCases := []responseTestCase{
		{
			output: monad{
				Data: outputData{
					SomeField: "some data",
				},
			},
			expectedHTTPStatus:   200,
			expectedJSONResponse: []byte(`{"ok":{"data":"success: some data"}}`),
		},
		{
			output: monad{
				Error: errors.New("some output error"),
			},
			expectedHTTPStatus:   500,
			expectedJSONResponse: []byte(`{"error":{"message":"an error occurred: some output error"}}`),
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		rw := httptest.NewRecorder()
		presenter.Present(ctx, rw, testCase.output)
		response := rw.Result()

		if response.StatusCode != testCase.expectedHTTPStatus {
			t.Fatalf("wrong status code, expected %d, got %d", testCase.expectedHTTPStatus, response.StatusCode)
			break
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("could not read the response body: %S", err)
			break
		}

		if string(body) != string(testCase.expectedJSONResponse) {
			t.Fatalf("unexpected data, expected %s, got %s", string(testCase.expectedJSONResponse), string(body))
		}
	}
}
