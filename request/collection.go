package request

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
)

// CollectionID header and cookie keys
const (
	CollectionIDHeaderKey = "Collection-Id"
	CollectionIDCookieKey = "collection"
)

var (
	ErrInvalidCollectionId = errors.New("invalid collection id")
)

// GetCollectionID gets the collection id from the request
func GetCollectionID(req *http.Request) (string, error) {
	var collectionID string

	id, err := headers.GetCollectionID(req)
	if err == nil {
		collectionID = id
	} else if headers.IsErrNotFound(err) {
		collectionID, err = getCollectionIDFromCookie(req)
	}

	if err != nil {
		return "", err
	}

	if collectionID != "" && !isCollectionIDValid(collectionID) {
		collectionID = ""
		err = ErrInvalidCollectionId
	}

	return collectionID, err
}

func getCollectionIDFromCookie(req *http.Request) (string, error) {
	var collectionID string
	var err error

	c, err := req.Cookie(CollectionIDCookieKey)

	switch err {
	case nil:
		collectionID = c.Value
	case http.ErrNoCookie:
		err = nil // we don't consider this scenario an error so we set err to nil and return an empty string
	}

	return collectionID, err
}

var collectionIDPattern = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

func isCollectionIDValid(collectionID string) bool {
	return collectionIDPattern.MatchString(collectionID)
}
