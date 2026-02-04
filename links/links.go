package links

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// PathTransformer defines how to transform a path before building the URL
type PathTransformer func(path string) string

// GenericBuilder constructs URLs with custom path transformation logic
type GenericBuilder struct {
	BaseURL         *url.URL
	PathTransformer PathTransformer
}

// NewGenericBuilder creates a GenericBuilder with a custom path transformer
func NewGenericBuilder(
	baseURL *url.URL,
	transformer PathTransformer,
) *GenericBuilder {
	if transformer == nil {
		transformer = func(p string) string { return p }
	}
	return &GenericBuilder{
		BaseURL:         baseURL,
		PathTransformer: transformer,
	}
}

// FromHeadersOrDefaultGeneric creates a GenericBuilder from headers with
// custom transformer
func FromHeadersOrDefaultGeneric(
	h *http.Header,
	defaultURL *url.URL,
	transformer PathTransformer,
) *GenericBuilder {
	path := h.Get("X-Forwarded-Path-Prefix")
	host := h.Get("X-Forwarded-Host")

	baseURL := defaultURL
	if strings.HasPrefix(host, "api.") {
		baseURL = &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   path,
		}
	} else if path != "" {
		baseURL = defaultURL.JoinPath(path)
	}

	return NewGenericBuilder(baseURL, transformer)
}

// BuildURL transforms the old URL path and rebuilds it
func (gb *GenericBuilder) BuildURL(oldURL *url.URL) *url.URL {
	transformedPath := gb.PathTransformer(oldURL.Path)
	newURL := gb.BaseURL.JoinPath(transformedPath)
	newURL.RawQuery = oldURL.RawQuery
	return newURL
}

// BuildLink parses a link string and builds a new URL
func (gb *GenericBuilder) BuildLink(link string) (string, error) {
	if strings.TrimSpace(link) == "" {
		oldURL := &url.URL{Path: ""}
		return gb.BuildURL(oldURL).String(), nil
	}

	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}
	return gb.BuildURL(oldURL).String(), nil
}

func RemovePrefixFromPath(path, prefix string) string {
	for strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
	}
	return path
}

/* #################
	Standard transformers
 ################### */

// NoOpTransformer performs no transformation
func NoOpTransformer(path string) string {
	return path
}

// V1PrefixTransformer removes /v1 prefix
func V1PrefixTransformer(path string) string {
	return RemovePrefixFromPath(path, "/v1")
}

// DownloadPrefixTransformer removes /downloads prefix
func DownloadPrefixTransformer(path string) string {
	newPath := RemovePrefixFromPath(path, "/downloads")
	return newPath
}

// DownloadFilesPrefixTransformer removes /downloads/files prefix
func DownloadFilesPrefixTransformer(path string) string {
	return RemovePrefixFromPath(path, "/downloads/files")
}

/* #################
	Convenience builders
 ################### */

// NewDownloadBuilder creates a builder for download links
func NewDownloadBuilder(baseURL *url.URL) *GenericBuilder {
	return NewGenericBuilder(
		baseURL.JoinPath("downloads"),
		DownloadPrefixTransformer,
	)
}

// NewDownloadFilesBuilder creates a builder for download files links
func NewDownloadFilesBuilder(baseURL *url.URL) *GenericBuilder {
	return NewGenericBuilder(
		baseURL.JoinPath("downloads/files"),
		DownloadFilesPrefixTransformer,
	)
}

/* #################
	Custom transformer
 ################### */

// CustomTransformer allows arbitrary path transformation logic
func CustomTransformer(transform func(string) string) PathTransformer {
	return transform
}

// Deprecated: Use NewGenericBuilder with appropriate transformer instead.
// Will be removed in dp-net v4.0.0
type Builder struct {
	URL *url.URL
}

// Deprecated: Use NewGenericBuilder with appropriate transformer instead.
// Will be removed in dp-net v4.0.0
//
// Example replacement:
//
//	builder := NewGenericBuilder(baseURL, V1PrefixTransformer)
func (b *Builder) BuildURL(oldURL *url.URL) *url.URL {
	newPath := RemovePrefixFromPath(oldURL.Path, "/v1")
	apiURL := b.URL.JoinPath(newPath)
	apiURL.RawQuery = oldURL.RawQuery
	return apiURL
}

// Deprecated: Use NewGenericBuilder with DownloadPrefixTransformer instead.
// Will be removed in dp-net v4.0.0
//
// Example replacement:
//
//	builder := NewGenericBuilder(defaultURL, DownloadPrefixTransformer)
//	link, err := builder.BuildLink(oldLink)
func BuildDownloadLink(link string, defaultURL *url.URL) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}

	newPath := RemovePrefixFromPath(oldURL.Path, "/downloads")
	apiURL := defaultURL.JoinPath("downloads", newPath)
	apiURL.RawQuery = oldURL.RawQuery

	return apiURL.String(), nil
}

// Deprecated: Use NewGenericBuilder with DownloadFilesPrefixTransformer
// instead. Will be removed in dp-net v4.0.0
//
// Example replacement:
//
//	builder := NewGenericBuilder(defaultURL, DownloadFilesPrefixTransformer)
//	link, err := builder.BuildLink(oldLink)
func BuildDownloadFilesLink(
	link string,
	defaultURL *url.URL,
) (string, error) {
	oldURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse link to URL")
	}

	newPath := RemovePrefixFromPath(oldURL.Path, "/downloads/files")
	apiURL := defaultURL.JoinPath("downloads/files", newPath)
	apiURL.RawQuery = oldURL.RawQuery

	return apiURL.String(), nil
}
