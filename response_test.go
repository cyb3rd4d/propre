package propre_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyb3rd4d/propre"
)

type payload[T any] struct {
	OK            *T              `json:"ok,omitempty"`
	Error         *errorViewModel `json:"error,omitempty"`
	encodingError bool
}

func (p payload[T]) ContentType() string {
	return "application/json"
}

func (p payload[T]) Encode() ([]byte, error) {
	if p.encodingError {
		return nil, errors.New("encoding error")
	}

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

type responseTestCase struct {
	response             *propre.HTTPResponse[payload[okViewModel]]
	payload              payload[okViewModel]
	expectedHTTPStatus   int
	expectedJSONResponse []byte
}

func TestResponse(t *testing.T) {
	internalErrorBody := `{"error":"serious internal error"}`

	withHTTPResponseHeaders := propre.WithHTTPResponseHeaders[payload[okViewModel]](http.Header{
		"content-encoding": []string{"plain"},
		"x-custom-header":  []string{"custom-header-value"},
	})

	responseWithCustomInternalError := propre.NewHTTPResponse(
		withHTTPResponseHeaders,
		propre.WithGenericInternalError[payload[okViewModel]]([]byte(internalErrorBody)),
	)

	responseWithoutCustomInternalError := propre.NewHTTPResponse(
		withHTTPResponseHeaders,
	)

	testCases := map[string]responseTestCase{
		"success response": {
			response: responseWithCustomInternalError,
			payload: payload[okViewModel]{
				OK: &okViewModel{
					Data: "some data",
				},
			},
			expectedHTTPStatus:   200,
			expectedJSONResponse: []byte(`{"ok":{"data":"some data"}}`),
		},
		"error response": {
			response: responseWithCustomInternalError,
			payload: payload[okViewModel]{
				Error: &errorViewModel{
					Message: "some output error",
				},
			},
			expectedHTTPStatus:   500,
			expectedJSONResponse: []byte(`{"error":{"message":"some output error"}}`),
		},
		"response encoding error with custom internal error": {
			response: responseWithCustomInternalError,
			payload: payload[okViewModel]{
				encodingError: true,
			},
			expectedHTTPStatus:   500,
			expectedJSONResponse: []byte(internalErrorBody),
		},
		"response encoding error without custom error": {
			response: responseWithoutCustomInternalError,
			payload: payload[okViewModel]{
				encodingError: true,
			},
			expectedHTTPStatus:   500,
			expectedJSONResponse: []byte("internal error"),
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			rw := httptest.NewRecorder()
			testCase.response.Send(ctx, rw, testCase.payload)

			response := rw.Result()

			if response.StatusCode != testCase.expectedHTTPStatus {
				t.Fatalf("wrong status code, expected %d, got %d", testCase.expectedHTTPStatus, response.StatusCode)
				return
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("could not read the response body: %s", err)
				return
			}

			if string(body) != string(testCase.expectedJSONResponse) {
				t.Fatalf("unexpected data, expected %s, got %s", string(testCase.expectedJSONResponse), string(body))
				return
			}

			contentEncodingHeader := response.Header.Get("content-encoding")
			if contentEncodingHeader == "" {
				t.Fatal("content-encoding header not found in response")
				return
			}

			if contentEncodingHeader != "plain" {
				t.Fatalf("wrong content-encoding header value in response, expected %s, got %s", "plain", contentEncodingHeader)
				return
			}

			xCustomHeader := response.Header.Get("x-custom-header")
			if xCustomHeader == "" {
				t.Fatal("x-custom-header not found in response")
				return
			}

			if xCustomHeader != "custom-header-value" {
				t.Fatalf("wrong x-custom-header header value in response, expected %s, got %s", "custom-header-value", xCustomHeader)
				return
			}
		})
	}
}
