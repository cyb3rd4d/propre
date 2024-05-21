package propre

import "context"

// UseCaseHandler is the interface your use cases must implement.
// The purpose of an interactor (the implementation of a use case) is to deal with business rules.
// To do that an interactor can take domain components as dependencies such as repositories.
//
// The input can contain errors raised by the request decoder, the first thing to do is to check for
// that errors and to return domain errors if any.
//
// Then your domain objects can be manipulated and the produced output can hold either the successful
// scenario with the data to return to the client, or an error. The output will then be handled by the
// response sender to produce the appropriate HTTP response.
type UseCaseHandler[Input any, Output any] interface {
	Handle(context.Context, Input) Output
}
