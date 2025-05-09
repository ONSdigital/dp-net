package handlers

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"

	clientsidentity "github.com/ONSdigital/dp-api-clients-go/v2/identity"
	dphttp "github.com/ONSdigital/dp-net/v3/http"
	dprequest "github.com/ONSdigital/dp-net/v3/request"
	"github.com/ONSdigital/log.go/v2/log"
)

type getTokenFromReqFunc func(ctx context.Context, r *http.Request) (string, error)

// Identity is a handler that controls the authenticating of a request
func Identity(zebedeeURL string) func(http.Handler) http.Handler {
	authClient := clientsidentity.New(zebedeeURL)
	return IdentityWithHTTPClient(authClient)
}

// IdentityWithHTTPClient allows a handler to be created that uses the given identity client
func IdentityWithHTTPClient(cli *clientsidentity.Client) func(http.Handler) http.Handler {
	// maintain the public interface to ensure backwards compatible and allow the get X token functions to be passed into the handler func.
	return identityWithHTTPClient(cli, GetFlorenceToken, getServiceAuthToken)
}

func identityWithHTTPClient(cli *clientsidentity.Client, getFlorenceToken, getServiceToken getTokenFromReqFunc) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			log.Info(ctx, "executing identity check middleware")

			florenceToken, err := getFlorenceToken(ctx, req)
			if err != nil {
				handleFailedRequest(ctx, w, req, http.StatusInternalServerError, "error getting florence access token from request", err, nil)
				return
			}

			serviceAuthToken, err := getServiceToken(ctx, req)
			if err != nil {
				handleFailedRequest(ctx, w, req, http.StatusInternalServerError, "error getting service access token from request", err, nil)
				return
			}

			ctx, statusCode, authFailure, err := cli.CheckRequest(req, florenceToken, serviceAuthToken)
			logData := log.Data{"auth_status_code": statusCode}

			if err != nil {
				handleFailedRequest(ctx, w, req, statusCode, "identity client check request returned an error", err, logData)
				return
			}

			if authFailure != nil {
				handleFailedRequest(ctx, w, req, statusCode, "identity client check request returned an auth error", authFailure, logData)
				log.Error(ctx, "identity client check request returned an auth error", authFailure, logData)
				return
			}

			log.Info(ctx, "identity client check request completed successfully invoking downstream http handler")

			req = req.WithContext(ctx)
			h.ServeHTTP(w, req)
		})
	}
}

// handleFailedRequest adhering to the DRY principle - clean up for failed identity requests, log the error, drain the request body and write the status code.
func handleFailedRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, event string, err error, data log.Data) {
	log.Error(ctx, event, err, data)
	dphttp.DrainBody(r)
	w.WriteHeader(status)
}

func GetFlorenceToken(ctx context.Context, req *http.Request) (string, error) {
	var florenceToken string

	token, err := headers.GetUserAuthToken(req)
	if err == nil {
		florenceToken = token
	} else if headers.IsErrNotFound(err) {
		log.Info(ctx, "florence access token header not found attempting to find access token cookie")
		florenceToken, err = getFlorenceTokenFromCookie(ctx, req)
	}

	return florenceToken, err
}

func getFlorenceTokenFromCookie(ctx context.Context, req *http.Request) (string, error) {
	var florenceToken string
	var err error

	c, err := req.Cookie(dprequest.FlorenceCookieKey)
	if err == nil {
		florenceToken = c.Value
	} else if err == http.ErrNoCookie {
		err = nil // we don't consider this scenario an error so we set err to nil and return an empty token
		log.Info(ctx, "florence access token cookie not found in request")
	}

	return florenceToken, err
}

func getServiceAuthToken(ctx context.Context, req *http.Request) (string, error) {
	var authToken string

	token, err := headers.GetServiceAuthToken(req)
	if err == nil {
		authToken = token
	} else if headers.IsErrNotFound(err) {
		err = nil // we don't consider this scenario an error so we set err to nil and return an empty token
		log.Info(ctx, "service auth token request header is not found")
	}

	return authToken, err
}
