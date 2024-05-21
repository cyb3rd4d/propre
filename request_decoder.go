package propre

import "net/http"

type RequestDecoder[Input any] interface {
	Decode(req *http.Request) Input
}
