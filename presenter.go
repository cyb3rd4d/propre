package propre

import (
	"context"
	"io"
)

// Presenter is the interface responsible of producing the correct response
// depending on the given output.
type Presenter[Output any, Writer io.Writer] interface {
	Present(context.Context, Writer, Output)
}
