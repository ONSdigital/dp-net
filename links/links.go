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
	if host == "" || !strings.HasPrefix(host, "api.") {
		return &Builder{
			URL: defaultURL.JoinPath(path),
		}
	}

	scheme := h.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "https"
	}

	return &Builder{
		URL: &url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   path,
		},
	}
}

func (b *Builder) BuildURL(oldURL *url.URL) *url.URL {
	newPath := oldURL.Path
	for strings.HasPrefix(newPath, "/v1") {
		newPath = strings.TrimPrefix(newPath, "/v1")
	}

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
