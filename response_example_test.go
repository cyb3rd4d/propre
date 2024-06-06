package propre_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cyb3rd4d/propre"
)

type successViewModel struct {
	Data string `json:"data"`
}

func (s successViewModel) ContentType(ctx context.Context) string {
	return "application/json"
}

func (s successViewModel) Encode(ctx context.Context) ([]byte, error) {
	return json.Marshal(s)
}

func (s successViewModel) StatusCode(ctx context.Context) int {
	return http.StatusOK
}

type invalidViewModel struct{}

func (s invalidPayload) ContentType(ctx context.Context) string {
	return "application/json"
}

func (s invalidPayload) Encode(ctx context.Context) ([]byte, error) {
	return nil, errors.New("some encoding error")
}

func (s invalidPayload) StatusCode(ctx context.Context) int {
	// The actual HTTP code will be 500
	return 42
}

// This example illustrates a valid scenario. A presenter builds a view model
// returning a 200 status code, a Content-Type header set to "application/json",
// and a response body with a `data` JSON object.
//
// The response will also contains two headers:
// * content-encoding with the value "bzip"
// * x-custom-header with the value "custom header value"
func ExampleHTTPResponse_success() {
	ctx := context.Background()

	response := propre.NewHTTPResponse(
		propre.WithHTTPResponseHeaders[successViewModel](http.Header{
			"content-encoding": []string{"bzip"},
			"x-custom-header":  []string{"custom header value"},
		}),
	)

	successPayload := successViewModel{
		Data: "success response payload",
	}

	rw := httptest.NewRecorder()
	response.Send(ctx, rw, successPayload)
	result := rw.Result()

	contentEncodingHeader := result.Header.Get("content-encoding")
	fmt.Println(contentEncodingHeader)

	xCustomHeader := result.Header.Get("x-custom-header")
	fmt.Println(xCustomHeader)

	var responseBody successViewModel
	json.NewDecoder(result.Body).Decode(&responseBody)

	fmt.Println(responseBody.Data)
	fmt.Println(result.StatusCode)

	// Output:
	// bzip
	// custom header value
	// success response payload
	// 200
}

// In this example the view model returns an encoding error. Because a
// custom internal error is set in the response, it is returned to the
// client with a 500 status code.
func ExampleHTTPResponse_custom_error() {
	ctx := context.Background()

	response := propre.NewHTTPResponse(
		propre.WithGenericInternalError[invalidPayload](
			[]byte(`{"error":"custom internal error"}`),
		),
	)

	invalidPayload := invalidPayload{}

	rw := httptest.NewRecorder()
	response.Send(ctx, rw, invalidPayload)
	result := rw.Result()

	var responseBody map[string]any
	json.NewDecoder(result.Body).Decode(&responseBody)

	fmt.Println(responseBody["error"])
	fmt.Println(result.StatusCode)

	// Output:
	// custom internal error
	// 500
}
