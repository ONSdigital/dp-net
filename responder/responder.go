package responder

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

// Responder is responsible for responding to http requests, providing methods for responding
// in with JSON and handling errors
type Responder struct{}

// New returns a new responder
func New() *Responder {
	return &Responder{}
}

// JSON responds to a HTTP request, expecting the response body
// to be marshall-able into JSON
func (r *Responder) JSON(ctx context.Context, w http.ResponseWriter, status int, resp interface{}){
	b, err := json.Marshal(resp)
	if err != nil {
		respondError(ctx, w, Error{
			statusCode: http.StatusInternalServerError,
			err:        errors.Wrap(err, "failed to marshal response"),
			message:    "Internal Server Error: Badly formed reponse attempt",
			logData: log.Data{
				"response": resp,
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err = w.Write(b); err != nil {
		log.Error(ctx, "failed to write response", err, log.Data{
			"response": string(b),
		})
		return
	}
}

// Error responds with a single error, formatted to fit in ONS's desired error
// response structure (essentially an array of errors)
func (r *Responder) Error(ctx context.Context, w http.ResponseWriter, err error) {
	respondError(ctx, w, err)
}

// respondError is the implementation of Error, seperated so it can be used internally
// by the other respond functions without having to create a new Responder
func respondError(ctx context.Context, w http.ResponseWriter, err error){
	log.Info(ctx, "error responding to HTTP request", log.ERROR, &log.EventErrors{{
			Message:    err.Error(),
			StackTrace: stackTrace(err),
			Data:       unwrapLogData(err),
		}},
	)

	status := unwrapStatusCode(err)
	msg    := errorMessage(err)
	resp   := errorResponse{
		Errors: []string{msg},
	}

	logData := log.Data{
		"error":       err.Error(),
		"response":    msg,
		"status_code": status,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Error(ctx, "badly formed error response", err, logData)
		http.Error(w, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "failed to write error response", err, logData)
		return
	}

	log.Info(ctx, "returned error response", logData)
}

// Errors responds with a slice of errors formatted to ONS's desired error response structure.
// Note you will have to pass a slice of []error rather than slice of []{some type that implements
// error}. The underlying structs can be any type that implements error but the slice itself must be
// defined as []error
func (r *Responder) Errors(ctx context.Context, w http.ResponseWriter, status int, errs []error){
	var errorLogs log.EventErrors
	var errorMsgs []string

	for _, err := range errs{
		errorLogs = append(errorLogs,log.EventError{
			Message:    err.Error(),
			StackTrace: stackTrace(err),
			Data:       unwrapLogData(err),
		})
		errorMsgs = append(errorMsgs, errorMessage(err))
	}

	log.Info(ctx, "error responding to HTTP request", log.ERROR, &errorLogs)

	resp   := errorResponse{
		Errors: errorMsgs,
	}

	logData := log.Data{
		"response":    resp,
		"status_code": status,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Error(ctx, "badly formed error response", err, logData)
		http.Error(w, "unexpected error", status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "failed to write error response", err, logData)
		return
	}

	log.Info(ctx, "returned error response", logData)
}

// Bytes responds to a http request with the raw bytes of whatever's passed as
// resp. Can be used to respond with a raw string, bytes, pre-encoded object etc
func (r *Responder) Bytes(ctx context.Context, w http.ResponseWriter, status int, resp []byte){
	w.WriteHeader(status)
	if _, err := w.Write(resp); err != nil {
		log.Error(ctx, "failed to write response", err, log.Data{
			"response": string(resp),
		})
		return
	}
}

// StatusCode responds with a raw status code
func (r *Responder) StatusCode(w http.ResponseWriter, status int){
	w.WriteHeader(status)
}