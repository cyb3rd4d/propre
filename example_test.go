package propre_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cyb3rd4d/propre"
)

type CreateTodoInput struct {
	Data struct {
		Title string
	}
	Error error
}

type CreateTodoOutput struct {
	Data struct {
		ID    int
		Title string
	}
	Error error
}

type CreateTodoRequestDecoder[Input CreateTodoInput] struct{}

func (decoder *CreateTodoRequestDecoder[Input]) Decode(req *http.Request) CreateTodoInput {
	var input CreateTodoInput
	var requestBody map[string]any
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		input.Error = fmt.Errorf("[CreateTodoRequestDecoder] %w caused by %s", errors.New("request read error"), err)
		return input
	}

	todoTitle, ok := requestBody["title"].(string)
	if !ok {
		input.Error = fmt.Errorf("[CreateTodoRequestDecoder] %w caused by %s", errors.New("missing title field"), err)
		return input
	}

	input.Data.Title = todoTitle
	return input
}

type CreateTodoUseCaseInteractor[Input CreateTodoInput, Output CreateTodoOutput] struct{}

func (g *CreateTodoUseCaseInteractor[Input, Output]) Handle(ctx context.Context, input CreateTodoInput) CreateTodoOutput {
	var output CreateTodoOutput
	if input.Error != nil {
		output.Error = fmt.Errorf("[CreateTodoUseCase] %w caused by %s", errors.New("input error"), input.Error)
		return output
	}

	output.Data.ID = 42
	output.Data.Title = input.Data.Title
	return output
}

type CreateTodoResponseSender[Output CreateTodoOutput] struct{}

func (s *CreateTodoResponseSender[Output]) Send(ctx context.Context, rw http.ResponseWriter, output CreateTodoOutput) {
	type responseBodyOK struct {
		Data struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
	}

	type responseBodyError struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if output.Error != nil {
		// Here you should use some kind of error mapping to return an appropriate response depending on
		// the error.
		var body responseBodyError
		body.Error.Message = "create todo error"
		s.sendResponse(ctx, rw, http.StatusInternalServerError, body)
		return
	}

	var body responseBodyOK
	body.Data.ID = output.Data.ID
	body.Data.Title = output.Data.Title
	s.sendResponse(ctx, rw, http.StatusCreated, body)
}

func (s *CreateTodoResponseSender[Output]) sendResponse(_ context.Context, rw http.ResponseWriter, statusCode int, body any) {
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(statusCode)
	err := json.NewEncoder(rw).Encode(body)
	if err != nil {
		// handle the error
	}
}

// In this example an HTTP handler is created to handle POST requests.
// The use case is the creation of a new task in a todo list.
//
// First, a request decoder read the request to extract the data from its body. Those data are stored
// in a `CreateTodoInput` struct and are passed to a use case handler.
//
// Then the use case handler checks if the input contains an error. If it is the case, it returns an output
// containing a business error corresponding to the request validation error. If the input does not contain
// an error, the task can be saved in the DB through a repository for example, and an output is returned
// containing the title of the task and its ID.
//
// Finally, it's time to return a response to the client. The output is passed to the response sender.
// If the output contains an error, an appropriate response must be returned. For that you can use a
// mapping, or a dedicated component. If the output does not contain an error, the response can be built
// with the data held in the struct.
func ExampleHTTPHandler() {
	requestDecoder := &CreateTodoRequestDecoder[CreateTodoInput]{}
	useCaseHandler := &CreateTodoUseCaseInteractor[CreateTodoInput, CreateTodoOutput]{}
	responseSender := &CreateTodoResponseSender[CreateTodoOutput]{}
	httpHandler := propre.NewHTTPHandler(requestDecoder, useCaseHandler, responseSender)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title":"New todo title"}`))
	httpHandler.ServeHTTP(rw, req)

	var data map[string]any
	err := json.NewDecoder(rw.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Println(rw.Code)
	fmt.Println(data)
	// Output:
	// 201
	// map[data:map[id:42 title:New todo title]]
}
