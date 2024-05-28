package propre_test

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/cyb3rd4d/propre"
)

type SomeRequestPayload struct {
	SomeField string `json:"some_field"`
}

func (p SomeRequestPayload) Validate() error {
	var err error
	if p.SomeField == "" {
		err = errors.New("empty field value")
	}

	return err
}

// In this example a request with a body considered as valid by the method
// SomeRequestPayload.Validate is passed to a JSON extractor.
// The resulting payload is typed as SomeRequestPayload and its field
// SomeField is accessible thanks to the generic type parameter.
func ExampleRequestPayloadExtractor_valid() {
	extractor := propre.NewRequestPayloadExtractor[SomeRequestPayload](
		propre.JSONDecoder,
	)

	req := httptest.NewRequest("GET", "/", io.NopCloser(strings.NewReader(
		`{"some_field":"some data"}`,
	)))

	payload, err := extractor.Extract(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(payload.SomeField)
	// Output: some data
}

// Here a request with a body considered as invalid by the method
// SomeRequestPayload.Validate is passed to a JSON extractor.
// The Extract method returns the error raised by the method
// SomeRequestPayload.Validate.
func ExampleRequestPayloadExtractor_invalid() {
	extractor := propre.NewRequestPayloadExtractor[SomeRequestPayload](
		propre.JSONDecoder,
	)

	req := httptest.NewRequest("GET", "/", io.NopCloser(strings.NewReader(
		`{"some_field":""}`,
	)))

	_, err := extractor.Extract(req)
	if err == nil {
		panic("should not be nil")
	}

	fmt.Println(err.Error())
	// Output: empty field value
}

// If a request body is malformed, the decoder will raise an error.
// This error is wrapped by ErrRequestPayloadExtraction and is returned.
func ExampleRequestPayloadExtractor_malformed_body() {
	extractor := propre.NewRequestPayloadExtractor[SomeRequestPayload](
		propre.JSONDecoder,
	)

	req := httptest.NewRequest("GET", "/", io.NopCloser(strings.NewReader(
		`{"some_field`,
	)))

	_, err := extractor.Extract(req)
	if err == nil {
		panic("should not be nil")
	}

	// This error is raised if the decoder failed to extract the request body.
	if errors.Is(err, propre.ErrRequestPayloadExtraction) {
		fmt.Println("decoder error")
	}

	// Output: decoder error
}
