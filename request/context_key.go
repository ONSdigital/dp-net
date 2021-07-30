package request

import v1request "github.com/ONSdigital/dp-net/request"

// Context keys - using type defined in v1 to prevent any type mismatch issue.
const (
	UserIdentityKey        = v1request.ContextKey("User-Identity")
	CallerIdentityKey      = v1request.ContextKey("Caller-Identity")
	RequestIdKey           = v1request.ContextKey("request-id")
	FlorenceIdentityKey    = v1request.ContextKey("florence-id")
	LocaleContextKey       = v1request.ContextKey(LocaleHeaderKey)
	CollectionIDContextKey = v1request.ContextKey(CollectionIDHeaderKey)
)
