package request

type ContextKey string

// Context keys - using type defined in v1 to prevent any type mismatch issue.
const (
	UserIdentityKey        = ContextKey("User-Identity")
	CallerIdentityKey      = ContextKey("Caller-Identity")
	RequestIdKey           = ContextKey("request-id")
	FlorenceIdentityKey    = ContextKey("florence-id")
	LocaleContextKey       = ContextKey(LocaleHeaderKey)
	CollectionIDContextKey = ContextKey(CollectionIDHeaderKey)
)
