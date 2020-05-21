package handlers

import dpHTTP "github.com/ONSdigital/dp-net/http"

// Key - iota enum of possible sets of keys for middleware manipulation
type Key int

// KeyMap represents a mapping between header keys, cookie keys and context keys for an equivalent value.
type KeyMap struct {
	Header  string
	Cookie  string
	Context dpHTTP.ContextKey
}

// Possible values for sets of keys
const (
	UserAccess Key = iota
	Locale
	CollectionID
)

// Header returns the header key
func (k Key) Header() string {
	return KeyMaps[k].Header
}

// Cookie returns the cookie key
func (k Key) Cookie() string {
	return KeyMaps[k].Cookie
}

// Context returns the context key
func (k Key) Context() dpHTTP.ContextKey {
	return KeyMaps[k].Context
}

// KeyMaps maps the possible values of Key enumeration to their Header, Cookie and Context correspnding keys
var KeyMaps = map[Key]*KeyMap{
	UserAccess: {
		Header:  dpHTTP.FlorenceHeaderKey,
		Cookie:  dpHTTP.FlorenceCookieKey,
		Context: dpHTTP.FlorenceIdentityKey,
	},
	Locale: {
		Header:  dpHTTP.LocaleHeaderKey,
		Cookie:  dpHTTP.LocaleCookieKey,
		Context: dpHTTP.ContextKey(dpHTTP.LocaleHeaderKey),
	},
	CollectionID: {
		Header:  dpHTTP.CollectionIDHeaderKey,
		Cookie:  dpHTTP.CollectionIDCookieKey,
		Context: dpHTTP.ContextKey(dpHTTP.CollectionIDHeaderKey),
	},
}
