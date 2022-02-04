package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-net/v2/responder"

	"github.com/pkg/errors"
)

type testRequest struct{
	Hello string `json:"hello"`
}

type testResponse struct{
	Message string `json:"message"`
}

type testHandler struct{
	respond responder.Responder
}

func (h *testHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// First error response takes full advantage of interfaces
	// asserted for. Logs original error string with attached log data
	// and stack trace (by wrapping with pkg/errors), responds to
	// user with specified message and status code.
	var req testRequest
	if err json.NewDecoder(r.Body).Decode(&req); err != nil{
		h.respond.Error(ctx, w, h.respond.Error(ctx, w, &testError{
			err:        errors.Wrap(err, "failed to decode"),
			statusCode: http.StatusBadRequest,
			logData:    log.Data{
				"log": "me",
			}
			message:    "bad formed request",
		)
		return
	}

	// Basic Go errors work too, defaults to status 500 and returns
	// original error string to user, no stack trace.
	if req.Hello != "world"{
		h.respond.Error(ctx, w, errors.New("Hello, world!"))
		return
	}

	// By making sure errors are wrapped at each level, logData, status codes
	// response messages and stack traces can be propagated down the call stack.
	// Status codes and error messages can be overwritten at any level,
	// stack trace always points to deepest point it was wrapped.
	if err := h.someFunc(); err != nil{
		h.respond.Error(ctx, w, fmt.Errorf("failed to someFunc: %w", err))
		return
	}

	resp := testResponse{
		Message: "hello",
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
})

func (h *testHandler) someFunc() error{
	// Any combination of information can be included in errors,
	// ommitting what's unnecessary.
	return &testError{
		err:        errors.New("original cause"),
		statusCode: http.StatusForbidden,
	}
}

func main() {
	handler := testHandler{
		respond: responder.New(),
	}

	mux := http.NewServeMux
	mux.HandleFunc("/hello", testHandler)

	http.ListenAndServe(":3333", mux)
}
