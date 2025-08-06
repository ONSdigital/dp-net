package links

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Builder struct {
	URL *url.URL
}

func FromHeadersOrDefault(h *http.Header, defaultURL *url.URL) *Builder {
	path := h.Get("X-Forwarded-Path-Prefix")

	host := h.Get("X-Forwarded-Host")
	if strings.HasPrefix(host, "api.") {
		return &Builder{
			URL: &url.URL{
				Scheme: "https",
				Host:   host,
				Path:   path,
			},
		}
	}

	return &Builder{
		URL: defaultURL.JoinPath(path),
	}
}

func (b *Builder) BuildURL(oldURL *url.URL) *url.URL {
	newPath := oldURL.Path
	newPath = RemovePrefixFromPath(newPath, "/v1")

	apiURL := b.URL.JoinPath(newPath)
	apiURL.RawQuery = oldURL.RawQuery
	return apiURL
}

func (b *Builder) BuildLink(link string) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}

	return b.BuildURL(oldURL).String(), nil
}

func BuildDownloadLink(link string, defaultURL *url.URL) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}

	newPath := oldURL.Path
	newPath = RemovePrefixFromPath(newPath, "/downloads")

	apiURL := defaultURL.JoinPath("downloads", newPath)
	apiURL.RawQuery = oldURL.RawQuery

	return apiURL.String(), nil
}

func BuildDownloadNewLink(link string, defaultURL *url.URL) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}

	newPath := oldURL.Path
	newPath = RemovePrefixFromPath(newPath, "/downloads-new")

	apiURL := defaultURL.JoinPath("downloads-new", newPath)
	apiURL.RawQuery = oldURL.RawQuery

	return apiURL.String(), nil
}

func RemovePrefixFromPath(path, prefix string) string {
	for strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
	}
	return path
}
