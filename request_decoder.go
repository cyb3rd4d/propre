package propre

import "net/http"

// The RequestDecoder's purpose is to check and extract the request's data required by a
// use case. The produced input can be either an error (the request does not match
// the requirements) or the actual data required by the use case.
type RequestDecoder[Input any] interface {
	Decode(req *http.Request) Input
}
