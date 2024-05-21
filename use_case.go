package propre

import "context"

type UseCaseHandler[Input any, Output any] interface {
	Handle(context.Context, Input) Output
}
