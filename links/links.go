package links

import (
	"context"
	"net/url"
)

type contextKey string

const (
	ctxAPIURL = "apiurl"
)

func URLBuild(ctx context.Context, oldURL url.URL) (string, error) {
	apiURL := ctx.Value(ctxAPIURL).(url.URL)
	apiURL.JoinPath(oldURL.Path)
	apiURL.RawQuery = oldURL.RawQuery
	return apiURL.String(), nil
}
