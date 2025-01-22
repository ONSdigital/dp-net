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
	h.Del("Authorization")
	h.Del("X-Florence-Token")

	r.Header.Del("Authorization")
	r.Header.Del("X-Florence-Token")

	log.Info(r.Context(), "building URL from headers", log.Data{
		"headers":        h,
		"requestHeaders": r.Header,
		"defaultURL":     defaultURL.String(),
		"url":            r.URL.String(),
	})

	path := h.Get("X-Forwarded-Path-Prefix")

	host := h.Get("X-Forwarded-Host")
	if host == "" || r.Host == "" {
		defaultURL = defaultURL.JoinPath(path)
		log.Info(r.Context(), "X-Forwarded-Host or r.Host is empty, using default URL", log.Data{
			"headers":        h,
			"requestHeaders": r.Header,
			"r.Host":         r.Host,
			"r.remoteAddr":   r.RemoteAddr,
		})
		return &Builder{
			URL: defaultURL,
		}
	}
	if !strings.HasPrefix(host, "api") {
		log.Info(r.Context(), "X-Forwarded-Host is not an external host, using incoming request host", log.Data{
			"headers":        h,
			"requestHeaders": r.Header,
			"r.Host":         r.Host,
			"r.remoteAddr":   r.RemoteAddr,
		})
		host = r.Host
	}

	scheme := h.Get("X-Forwarded-Proto")
	if scheme == "" {
		log.Info(r.Context(), "X-Forwarded-Proto is empty, using http or https based on host", log.Data{
			"headers":        h,
			"requestHeaders": r.Header,
			"host":           host,
		})
		if !strings.HasPrefix(host, "api") {
			scheme = "http"
		} else {
			scheme = "https"
		}
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
