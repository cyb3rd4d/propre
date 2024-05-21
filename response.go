package propre

import (
	"context"
	"net/http"
)

type HTTPResponseSender[Output any] interface {
	Send(context.Context, http.ResponseWriter, Output)
}
