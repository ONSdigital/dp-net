package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-net/v2/responder"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/pkg/errors"
)

type testRequest struct {
	Hello string `json:"hello"`
}

type testResponse struct {
	Message string `json:"message"`
}

type testHandler struct {
	respond *responder.Responder
}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// First error response takes full advantage of interfaces
	// asserted for. Logs original error string with attached log data
	// and stack trace (by wrapping with pkg/errors), responds to
	// user with specified message and status code.
	var req testRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &testError{
			err: errors.Wrap(err, "failed to decode"),
			logData: log.Data{
				"log": "me",
			},
			message: "badly formed request",
		})
		return
	}

	// Basic Go errors work too, logs and returns
	// original error string to user, no stack trace or log data.
	if req.Hello != "world" {
		h.respond.Error(ctx, w, http.StatusInternalServerError, errors.New("Hello, world!"))
		return
	}

	// By making sure errors are wrapped at each level, logData,
	// response messages and stack traces can be propagated down the call stack.
	// Status codes and error messages can be overwritten at any level,
	// stack trace always points to deepest point it was wrapped.
	if err := h.someFunc(); err != nil {
		h.respond.Error(ctx, w, http.StatusUnauthorized, fmt.Errorf("failed to someFunc: %w", err))
		return
	}

	resp := testResponse{
		Message: "hello",
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (h *testHandler) someFunc() error {
	// Any combination of information can be included in errors,
	// ommitting what's unnecessary.
	return &testError{
		err:     errors.New("original cause"),
		message: "message returned to user",
	}
}

func main() {
	handler := testHandler{
		respond: responder.New(),
	}

	mux := http.NewServeMux()
	mux.Handle("/hello", handler)

	panic(http.ListenAndServe(":3333", mux))
}

type testError struct {
	err        error
	statusCode int
	logData    map[string]interface{}
	message    string
}

// standard Go error interfaces
func (e *testError) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

func (e *testError) Unwrap() error {
	return e.err
}

func (e *testError) LogData() map[string]interface{} {
	return e.logData
}

func (e *testError) Message() string {
	return e.message
}
