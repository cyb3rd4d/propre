package propre_test

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyb3rd4d/propre"
)

var (
	errInvalidPayload = errors.New("invalid payload")
)

type validPayload struct {
	XMLName   xml.Name `json:"-" xml:"ValidPayload"`
	SomeField string   `json:"some_field" xml:"SomeField"`
}

func (p validPayload) Validate() error {
	return nil
}

type invalidPayload struct {
	XMLName   xml.Name `json:"-" xml:"InvalidPayload"`
	SomeField string   `json:"some_field" xml:"SomeField"`
}

func (p invalidPayload) Validate() error {
	return errInvalidPayload
}

type invalidPayloadDecoder struct{}

func (d *invalidPayloadDecoder) Decode(any) error {
	return errors.New("foo")
}

type testCase[T propre.Validatable] struct {
	req            *http.Request
	sut            func(req *http.Request) (T, error)
	isPayloadValid func(gotPayload any) bool
	expectedError  error
}

func TestExtractRequestPayload(t *testing.T) {
	testCases := map[string]testCase[propre.Validatable]{
		"valid JSON request": {
			req: httptest.NewRequest("GET", "/", io.NopCloser(
				strings.NewReader(`{"some_field":"some_data"}`),
			)),
			sut: func(req *http.Request) (propre.Validatable, error) {
				extractor := propre.NewRequestPayloadExtractor[validPayload](propre.JSONDecoder)
				return extractor.Extract(req)
			},
			isPayloadValid: func(gotPayload any) bool {
				return gotPayload.(validPayload).SomeField == "some_data"
			},
			expectedError: nil,
		},
		"valid XML request": {
			req: httptest.NewRequest("GET", "/", io.NopCloser(
				strings.NewReader(`<ValidPayload><SomeField>some_data</SomeField></ValidPayload>`),
			)),
			sut: func(req *http.Request) (propre.Validatable, error) {
				extractor := propre.NewRequestPayloadExtractor[validPayload](propre.XMLDecoder)
				return extractor.Extract(req)
			},
			isPayloadValid: func(gotPayload any) bool {
				return gotPayload.(validPayload).SomeField == "some_data"
			},
			expectedError: nil,
		},
		"invalid JSON request": {
			req: httptest.NewRequest("GET", "/", io.NopCloser(
				strings.NewReader(`{"some_field":"some_data"}`),
			)),
			sut: func(req *http.Request) (propre.Validatable, error) {
				extractor := propre.NewRequestPayloadExtractor[invalidPayload](propre.JSONDecoder)
				return extractor.Extract(req)
			},
			expectedError: errInvalidPayload,
		},
		"invalid XML request": {
			req: httptest.NewRequest("GET", "/", io.NopCloser(
				strings.NewReader(`<InvalidPayload><SomeField>some_data</SomeField></InvalidPayload>`),
			)),
			sut: func(req *http.Request) (propre.Validatable, error) {
				extractor := propre.NewRequestPayloadExtractor[invalidPayload](propre.XMLDecoder)
				return extractor.Extract(req)
			},
			expectedError: errInvalidPayload,
		},
		"decoder error": {
			req: httptest.NewRequest("GET", "/", nil),
			sut: func(req *http.Request) (propre.Validatable, error) {
				decoder := func(io.Reader) propre.PayloadDecoder {
					return &invalidPayloadDecoder{}
				}

				extractor := propre.NewRequestPayloadExtractor[invalidPayload](decoder)
				return extractor.Extract(req)
			},
			expectedError: propre.ErrRequestPayloadExtraction,
		},
	}

	for scenario, testCase := range testCases {
		t.Run(scenario, func(t *testing.T) {
			payload, err := testCase.sut(testCase.req)
			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("%s scenario failed: got nil error", scenario)
				}

				if errors.Is(err, testCase.expectedError) == false {
					t.Fatalf("%s scenario failed: unexpected error: %s", scenario, err.Error())
				}

				return
			}

			if testCase.isPayloadValid(payload) == false {
				t.Fatalf("%s scenario failed: unexpected payload: %#v", scenario, payload)
			}
		})
	}
}
