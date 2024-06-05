package propre

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// PayloadDecoder is required by [RequestPayloadExtractor] to be able to decode the
// body of a request.
//
// [JSONDecoder] and [XMLDecoder] are two functions that return standard decoders,
// but you can create your own for your specific needs.
type PayloadDecoder interface {
	Decode(payload any) error
}

// JSONDecoder can be used in [NewRequestPayloadExtractor] to decode JSON payloads.
func JSONDecoder(req io.Reader) PayloadDecoder {
	return json.NewDecoder(req)
}

// XMLDecoder can be used in [NewRequestPayloadExtractor] to decode XML payloads.
func XMLDecoder(req io.Reader) PayloadDecoder {
	return xml.NewDecoder(req)
}

// Validatable is the constraint required by [RequestPayloadExtractor].
// Once the payload is decoded from the the request body, the method Validate() is
// called. It is the responsibility of the payload struct to know how to validate
// itself.
type Validatable interface {
	Validate() error
}

// RequestPayloadExtractor is the component used to extract and validate a
// request body.
type RequestPayloadExtractor[Payload Validatable] struct {
	decoder func(io.Reader) PayloadDecoder
}

// NewRequestPayloadExtractor builds a new [RequestPayloadExtractor].
// It takes a decoder function as dependency to instantiate the correct decoder
// for the given Payload type.
//
// [JSONDecoder] and [XMLDecoder] are two functions that return standard decoders,
// but you can create your own for your specific needs.
func NewRequestPayloadExtractor[Payload Validatable](
	decoder func(io.Reader) PayloadDecoder,
) *RequestPayloadExtractor[Payload] {
	return &RequestPayloadExtractor[Payload]{
		decoder: decoder,
	}
}

// Extract takes a request as an argument and extracts its body into the given Payload type.
// If the decoder fails an error [ErrRequestPayloadExtraction] wraps the decoder error.
// Then the method Validate of the Payload type is called and its error is returned.
func (extractor *RequestPayloadExtractor[Payload]) Extract(req *http.Request) (Payload, error) {
	var payload Payload
	decoder := extractor.decoder(req.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, fmt.Errorf("%w caused by %s", ErrRequestPayloadExtraction, err)
	}

	return payload, payload.Validate()
}
