package links

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type Builder struct {
	URL *url.URL
}

func FromHeadersOrDefault(h *http.Header, defaultURL *url.URL) *Builder {
	host := h.Get("X-Forwarded-Host")
	if host == "" {
		return &Builder{
			URL: defaultURL,
		}
	}

	scheme := h.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}

	port := h.Get("X-Forwarded-Port")
	if port != "" {
		host += ":" + port
	}

	path := h.Get("X-Forwarded-Path-Prefix")

	url := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   "/",
	}

	return &Builder{
		URL: url.JoinPath(path),
	}
}

func (b *Builder) BuildURL(oldURL *url.URL) *url.URL {
	apiURL := *b.URL
	apiURL.JoinPath(oldURL.Path)
	apiURL.RawQuery = oldURL.RawQuery
	return &apiURL
}

func (b *Builder) BuildLink(link string) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unalble to parse link to URL")
	}
	newURL := b.BuildURL(oldURL)
	return newURL.String(), nil
}
