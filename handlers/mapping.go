package handlers

import (
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/go-ns/common"
)

// Key - iota enum of possible sets of keys for middleware manipulation
type Key int

// KeyMap represents a mapping between header keys, cookie keys and context keys for an equivalent value.
type KeyMap struct {
	Header  string
	Cookie  string
	Context common.ContextKey
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
func (k Key) Context() common.ContextKey {
	return KeyMaps[k].Context
}

// KeyMaps maps the possible values of Key enumeration to their Header, Cookie and Context correspnding keys
var KeyMaps = map[Key]*KeyMap{
	UserAccess: {
		Header:  dphttp.FlorenceHeaderKey,
		Cookie:  dphttp.FlorenceCookieKey,
		Context: common.FlorenceIdentityKey,
	},
	Locale: {
		Header:  dphttp.LocaleHeaderKey,
		Cookie:  dphttp.LocaleCookieKey,
		Context: common.ContextKey(dphttp.LocaleHeaderKey),
	},
	CollectionID: {
		Header:  dphttp.CollectionIDHeaderKey,
		Cookie:  dphttp.CollectionIDCookieKey,
		Context: common.ContextKey(dphttp.CollectionIDHeaderKey),
	},
}
