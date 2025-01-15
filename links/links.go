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

func FromHeadersOrDefault(h *http.Header, r *http.Request, defaultURL *url.URL) *Builder {
	path := h.Get("X-Forwarded-Path-Prefix")

	host := h.Get("X-Forwarded-Host")
	if host == "" {
		if r.Host != "" {
			defaultURL.Host = r.Host
		}
		defaultURL = defaultURL.JoinPath(path)
		return &Builder{
			URL: defaultURL,
		}
	}

	scheme := h.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "https"
	}

	port := h.Get("X-Forwarded-Port")
	if port != "" {
		host += ":" + port
	}

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
	newPath := oldURL.Path
	newPath = strings.TrimPrefix(newPath, "/v1")

	apiURL := b.URL.JoinPath(newPath)
	apiURL.RawQuery = oldURL.RawQuery
	return apiURL
}

func (b *Builder) BuildLink(link string) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}
	newURL := b.BuildURL(oldURL)
	return newURL.String(), nil
}
