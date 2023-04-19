package request

type ContextKey string

const (
	UserIdentityKey        = ContextKey("User-Identity")
	CallerIdentityKey      = ContextKey("Caller-Identity")
	RequestIdKey           = ContextKey("request-id")
	FlorenceIdentityKey    = ContextKey("florence-id")
	LocaleContextKey       = ContextKey(LocaleHeaderKey)
	CollectionIDContextKey = ContextKey(CollectionIDHeaderKey)
)
