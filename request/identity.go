package request

import (
	"context"
	"net/http"
)

// Header constants
const (
	AuthHeaderKey            = "Authorization"
	DownloadServiceHeaderKey = "X-Download-Service-Token"
	FlorenceHeaderKey        = "X-Florence-Token"
	UserHeaderKey            = "User-Identity"
)

// Cookie constants
const (
	FlorenceCookieKey = "access_token"
)

// Other identity-related constants
const (
	DeprecatedAuthHeader = "Internal-Token"
	LegacyUser           = "legacyUser"
	BearerPrefix         = "Bearer "
)

// CheckRequester is an interface to allow mocking of auth.CheckRequest
type CheckRequester interface {
	CheckRequest(*http.Request) (context.Context, int, error)
}

// IdentityResponse represents the response from the identity service
type IdentityResponse struct {
	Identifier string `json:"identifier"`
}

// IsUserPresent determines if a user identity is present on the given context
func IsUserPresent(ctx context.Context) bool {
	userIdentity := ctx.Value(UserIdentityKey)
	return userIdentity != nil && userIdentity != ""

}

// IsFlorenceIdentityPresent determines if a florence identity is present on the given context
func IsFlorenceIdentityPresent(ctx context.Context) bool {
	florenceID := ctx.Value(FlorenceIdentityKey)
	return florenceID != nil && florenceID != ""
}

// AddUserHeader sets the given user ID on the given request
func AddUserHeader(r *http.Request, user string) {
	r.Header.Add(UserHeaderKey, user)
}

// AddServiceTokenHeader sets the given service token on the given request
func AddServiceTokenHeader(r *http.Request, serviceToken string) {
	if len(serviceToken) > 0 {
		r.Header.Add(AuthHeaderKey, BearerPrefix+serviceToken)
	}
}

// AddDownloadServiceTokenHeader sets the given download service token on the given request
func AddDownloadServiceTokenHeader(r *http.Request, serviceToken string) {
	if len(serviceToken) > 0 {
		r.Header.Add(DownloadServiceHeaderKey, serviceToken)
	}
}

// User gets the user identity from the context
func User(ctx context.Context) string {
	userIdentity, _ := ctx.Value(UserIdentityKey).(string)
	return userIdentity
}

// SetUser sets the user identity on the context
func SetUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, UserIdentityKey, user)
}

// SetFlorenceIdentity sets the florence identity for authentication
func SetFlorenceIdentity(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, FlorenceIdentityKey, user)
}

// SetFlorenceHeader sets a florence Header if the corresponding Identity key is in context
func SetFlorenceHeader(ctx context.Context, r *http.Request) {
	if IsFlorenceIdentityPresent(ctx) {
		r.Header.Set(FlorenceHeaderKey, ctx.Value(FlorenceIdentityKey).(string))
	}
}

// AddFlorenceHeader sets the given user access token (florence token) token on the given request
func AddFlorenceHeader(r *http.Request, userAccessToken string) {
	if len(userAccessToken) > 0 {
		r.Header.Add(FlorenceHeaderKey, userAccessToken)
	}
}

// AddAuthHeaders sets authentication headers for request
func AddAuthHeaders(ctx context.Context, r *http.Request, serviceToken string) {
	if IsUserPresent(ctx) {
		AddUserHeader(r, User(ctx))
	}
	AddServiceTokenHeader(r, serviceToken)
}

// AddDeprecatedHeader sets the deprecated header on the given request
func AddDeprecatedHeader(r *http.Request, token string) {
	if len(token) > 0 {
		r.Header.Add(DeprecatedAuthHeader, token)
	}
}

// IsCallerPresent determines if an identity is present on the given context.
func IsCallerPresent(ctx context.Context) bool {
	callerIdentity := ctx.Value(CallerIdentityKey)
	isPresent := callerIdentity != nil && callerIdentity != ""

	return isPresent
}

// Caller gets the caller identity from the context
func Caller(ctx context.Context) string {
	callerIdentity, _ := ctx.Value(CallerIdentityKey).(string)
	return callerIdentity
}

// SetCaller sets the caller identity on the context
func SetCaller(ctx context.Context, caller string) context.Context {
	return context.WithValue(ctx, CallerIdentityKey, caller)
}
