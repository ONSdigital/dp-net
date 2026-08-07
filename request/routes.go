package request

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ErrorDescriptionMissingParameters = "unable to process request due to missing required route parameters"
)

var (
	ErrorMissingParameters = errors.New(ErrorDescriptionMissingParameters)
)

// CreateMissingRouteVariableError creates an error with a standardised message for if a route variable is missing
func createMissingRouteVariableError(variable string) error {
	return fmt.Errorf("missing route variable %s", variable)
}

// CreateMissingRouteVariableError creates an error with a standardised message for if a route variable is provided but empty
func createEmptyRouteVariableError(variable string) error {
	return fmt.Errorf("route variable %s cannot be empty", variable)
}

// GetRouteVariable gets the specified variable from the route variables, and returns the found value (if any), or an error if missing/empty
func GetRouteVariable(r *http.Request, variable string) (*string, error) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		return nil, ErrorMissingParameters
	}

	variable, ok := vars[variable]
	if !ok {
		return nil, createMissingRouteVariableError(variable)
	}

	if variable == "" {
		return nil, createEmptyRouteVariableError(variable)
	}

	return &variable, nil
}
