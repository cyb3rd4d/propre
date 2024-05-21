package propre

import (
	"context"
	"net/http"
)

// HTTPResponseSender is the interface responsible of producing the correct HTTP response
// depending on the given output.
type HTTPResponseSender[Output any] interface {
	Send(context.Context, http.ResponseWriter, Output)
}
