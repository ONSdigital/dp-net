package links

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

type Builder struct {
	URL *url.URL
}

func FromHeadersOrDefault(h *http.Header, r *http.Request, defaultURL *url.URL) *Builder {
	log.Info(r.Context(), "building URL from headers", log.Data{
		"headers":    h,
		"defaultURL": defaultURL.String(),
		"host":       r.Host,
		"url":        r.URL.String(),
	})
	path := h.Get("X-Forwarded-Path-Prefix")

	host := h.Get("X-Forwarded-Host")
	if host == "" {
		if r.Host != "" {
			defaultURL.Host = r.Host
		}
		defaultURL = defaultURL.JoinPath(path)
		log.Info(r.Context(), "internal request, using default URL", log.Data{
			"defaultURL": defaultURL.String(),
		})
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

	loggingURL := url.JoinPath(path)
	log.Info(r.Context(), "external request, using URL from headers", log.Data{
		"url": loggingURL.String(),
	})
	return &Builder{
		URL: url.JoinPath(path),
	}
}

func (b *Builder) BuildURL(oldURL *url.URL) *url.URL {
	newPath := oldURL.Path
	newPath = strings.TrimPrefix(newPath, "/v1")

	apiURL := b.URL.JoinPath(newPath)
	apiURL.RawQuery = oldURL.RawQuery
	log.Info(context.TODO(), "built URL", log.Data{
		"oldURL":  oldURL.String(),
		"newURL":  apiURL.String(),
		"builder": b.URL.String(),
	})
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
