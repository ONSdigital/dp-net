package links

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	defaultInternalURL = &url.URL{Scheme: "http", Host: "localhost:8080"}
	defaultDownloadURL = &url.URL{Scheme: "https", Host: "download.api.host"}
)

const (
	invalidURL = ":invalid/url"
)

// Test GenericBuilder with NoOpTransformer
func Test_GenericBuilder_NoOp(t *testing.T) {
	Convey("Given a GenericBuilder configured with NoOpTransformer",
		t, func() {
			Convey("When building links with various old URLs", func() {
				tests := []struct {
					description string
					builderURL  string
					oldLink     string
					expected    string
				}{
					{
						description: "an empty old link is provided",
						builderURL:  "http://localhost:8080",
						oldLink:     "",
						expected:    "http://localhost:8080",
					},
					{
						description: "an old link with no path is provided",
						builderURL:  "http://localhost:8080",
						oldLink:     "http://localhost:8080",
						expected:    "http://localhost:8080",
					},
					{
						description: "an old link with a different base URL " +
							"is provided",
						builderURL: "http://localhost:8080",
						oldLink:    "https://oldHost:1000",
						expected:   "http://localhost:8080",
					},
					{
						description: "an old link with a path is provided",
						builderURL:  "http://localhost:8080",
						oldLink:     "http://localhost:8080/some/path",
						expected:    "http://localhost:8080/some/path",
					},
					{
						description: "an old link with a path and different " +
							"base URL is provided",
						builderURL: "http://localhost:8080",
						oldLink:    "http://oldHost:1000/some/path",
						expected:   "http://localhost:8080/some/path",
					},
					{
						description: "an old link without a base URL is " +
							"provided",
						builderURL: "http://localhost:8080",
						oldLink:    "/some/path",
						expected:   "http://localhost:8080/some/path",
					},
					{
						description: "an old link with query parameters is " +
							"provided",
						builderURL: "http://localhost:8080",
						oldLink: "http://localhost:8080/some/path?" +
							"param1=value1&param2=value2",
						expected: "http://localhost:8080/some/path?" +
							"param1=value1&param2=value2",
					},
				}

				for _, tt := range tests {
					Convey(tt.description, func() {
						// Given
						bu, err := url.Parse(tt.builderURL)
						So(err, ShouldBeNil)
						builder := NewGenericBuilder(
							bu,
							NoOpTransformer,
						)

						// When
						newurl, err := builder.BuildLink(tt.oldLink)

						// Then
						So(err, ShouldBeNil)
						So(newurl, ShouldEqual, tt.expected)
						So(builder.BaseURL.String(),
							ShouldEqual, tt.builderURL)
					})
				}
			})

			Convey("When an invalid old URL is provided", func() {
				// Given
				builder := NewGenericBuilder(
					&url.URL{},
					NoOpTransformer,
				)

				// When
				newurl, err := builder.BuildLink(invalidURL)

				// Then
				So(err, ShouldNotBeNil)
				So(err.Error(),
					ShouldContainSubstring, "unable to parse link to URL")
				So(newurl, ShouldBeEmpty)
			})
		})
}

// Test GenericBuilder with V1PrefixTransformer
func Test_GenericBuilder_V1Prefix(t *testing.T) {
	Convey("Given a GenericBuilder configured with V1PrefixTransformer",
		t, func() {
			Convey("When building links that contain /v1 prefixes", func() {
				tests := []struct {
					description string
					builderURL  string
					oldLink     string
					expected    string
				}{
					{
						description: "an external link with /v1 prefix " +
							"is redirected to internal URL",
						builderURL: "http://localhost:8080",
						oldLink:    "https://some.api.host/v1/some/path",
						expected:   "http://localhost:8080/some/path",
					},
					{
						description: "an external link with /v1 prefix " +
							"is redirected to external URL",
						builderURL: "https://some.api.host/v1",
						oldLink:    "https://some.api.host/v1/some/path",
						expected:   "https://some.api.host/v1/some/path",
					},
					{
						description: "an external link with multiple /v1 " +
							"prefixes is normalized",
						builderURL: "https://some.api.host/v1",
						oldLink: "https://some.api.host/v1/v1/v1/" +
							"some/path",
						expected: "https://some.api.host/v1/some/path",
					},
					{
						description: "an internal link is redirected to " +
							"external URL with /v1 prefix",
						builderURL: "https://some.api.host/v1",
						oldLink:    "http://localhost:8080/some/path",
						expected:   "https://some.api.host/v1/some/path",
					},
					{
						description: "an internal link with query parameters " +
							"is redirected to external URL",
						builderURL: "https://some.api.host/v1",
						oldLink: "http://localhost:8080/some/path?" +
							"param1=value1&param2=value2",
						expected: "https://some.api.host/v1/some/path?" +
							"param1=value1&param2=value2",
					},
				}

				for _, tt := range tests {
					Convey(tt.description, func() {
						// Given
						bu, err := url.Parse(tt.builderURL)
						So(err, ShouldBeNil)
						builder := NewGenericBuilder(
							bu,
							V1PrefixTransformer,
						)

						// When
						newurl, err := builder.BuildLink(tt.oldLink)

						// Then
						So(err, ShouldBeNil)
						So(newurl, ShouldEqual, tt.expected)
					})
				}
			})
		})
}

// Test NewDownloadBuilder convenience function
func Test_NewDownloadBuilder(t *testing.T) {
	Convey("Given a NewDownloadBuilder configured with "+
		"DownloadPrefixTransformer",
		t, func() {
			Convey("When building download links", func() {
				tests := []struct {
					description string
					oldLink     string
					expected    string
				}{
					{
						description: "an empty old link is provided",
						oldLink:     "",
						expected: "https://download.api.host/" +
							"downloads",
					},
					{
						description: "an old link with only /downloads " +
							"prefix is provided",
						oldLink: "https://download.api.host/" +
							"downloads",
						expected: "https://download.api.host/" +
							"downloads",
					},
					{
						description: "an old link with /downloads and " +
							"path is provided",
						oldLink: "https://download.api.host/" +
							"downloads/some/path",
						expected: "https://download.api.host/" +
							"downloads/some/path",
					},
					{
						description: "an old link with no path is " +
							"provided",
						oldLink: "http://localhost:23600",
						expected: "https://download.api.host/" +
							"downloads",
					},
					{
						description: "an old link with a different base " +
							"URL is provided",
						oldLink: "https://localhost:23600",
						expected: "https://download.api.host/" +
							"downloads",
					},
					{
						description: "an old link with a path is " +
							"provided",
						oldLink: "http://localhost:23600/some/path",
						expected: "https://download.api.host/" +
							"downloads/some/path",
					},
					{
						description: "an old link without a base URL " +
							"is provided",
						oldLink: "/some/path",
						expected: "https://download.api.host/" +
							"downloads/some/path",
					},
					{
						description: "an old link without a base URL " +
							"and / prefix is provided",
						oldLink: "some/path",
						expected: "https://download.api.host/" +
							"downloads/some/path",
					},
					{
						description: "an old link with query parameters " +
							"is provided",
						oldLink: "http://localhost:23600/some/path?" +
							"param1=value1&param2=value2",
						expected: "https://download.api.host/" +
							"downloads/some/path?" +
							"param1=value1&param2=value2",
					},
					{
						description: "an old link with multiple " +
							"/downloads prefixes is provided",
						oldLink: "https://download.api.host/" +
							"downloads/downloads/downloads/some/path",
						expected: "https://download.api.host/" +
							"downloads/some/path",
					},
				}

				for _, tt := range tests {
					Convey(tt.description, func() {
						// Given
						builder := NewDownloadBuilder(defaultDownloadURL)

						// When
						newurl, err := builder.BuildLink(tt.oldLink)

						// Then
						So(err, ShouldBeNil)
						So(newurl, ShouldEqual, tt.expected)
					})
				}
			})

			Convey("When an invalid old URL is provided", func() {
				// Given
				builder := NewDownloadBuilder(defaultDownloadURL)

				// When
				newurl, err := builder.BuildLink(invalidURL)

				// Then
				So(err, ShouldNotBeNil)
				So(err.Error(),
					ShouldContainSubstring, "unable to parse link to URL")
				So(newurl, ShouldBeEmpty)
			})
		})
}

// Test NewDownloadFilesBuilder convenience function
func Test_NewDownloadFilesBuilder(t *testing.T) {
	Convey("Given a NewDownloadFilesBuilder configured with "+
		"DownloadFilesPrefixTransformer",
		t, func() {
			Convey("When building download files links", func() {
				tests := []struct {
					description string
					oldLink     string
					expected    string
				}{
					{
						description: "an empty old link is provided",
						oldLink:     "",
						expected: "https://download.api.host/" +
							"downloads/files",
					},
					{
						description: "an old link with only " +
							"/downloads/files prefix is provided",
						oldLink: "https://download.api.host/" +
							"downloads/files",
						expected: "https://download.api.host/" +
							"downloads/files",
					},
					{
						description: "an old link with /downloads/files " +
							"and path is provided",
						oldLink: "https://download.api.host/" +
							"downloads/files/some/path",
						expected: "https://download.api.host/" +
							"downloads/files/some/path",
					},
					{
						description: "an old link with no path is " +
							"provided",
						oldLink: "http://localhost:23600",
						expected: "https://download.api.host/" +
							"downloads/files",
					},
					{
						description: "an old link with a different base " +
							"URL is provided",
						oldLink: "https://localhost:23600",
						expected: "https://download.api.host/" +
							"downloads/files",
					},
					{
						description: "an old link with a path is " +
							"provided",
						oldLink: "http://localhost:23600/some/path",
						expected: "https://download.api.host/" +
							"downloads/files/some/path",
					},
					{
						description: "an old link without a base URL " +
							"is provided",
						oldLink: "/some/path",
						expected: "https://download.api.host/" +
							"downloads/files/some/path",
					},
					{
						description: "an old link without a base URL " +
							"and / prefix is provided",
						oldLink: "some/path",
						expected: "https://download.api.host/" +
							"downloads/files/some/path",
					},
					{
						description: "an old link with query parameters " +
							"is provided",
						oldLink: "http://localhost:23600/some/path?" +
							"param1=value1&param2=value2",
						expected: "https://download.api.host/" +
							"downloads/files/some/path?" +
							"param1=value1&param2=value2",
					},
					{
						description: "an old link with multiple " +
							"/downloads/files prefixes is provided",
						oldLink: "https://download.api.host/" +
							"downloads/files/downloads/files/" +
							"downloads/files/some/path",
						expected: "https://download.api.host/" +
							"downloads/files/some/path",
					},
				}

				for _, tt := range tests {
					Convey(tt.description, func() {
						// Given
						builder := NewDownloadFilesBuilder(
							defaultDownloadURL,
						)

						// When
						newurl, err := builder.BuildLink(tt.oldLink)

						// Then
						So(err, ShouldBeNil)
						So(newurl, ShouldEqual, tt.expected)
					})
				}
			})

			Convey("When an invalid old URL is provided", func() {
				// Given
				builder := NewDownloadFilesBuilder(defaultDownloadURL)

				// When
				newurl, err := builder.BuildLink(invalidURL)

				// Then
				So(err, ShouldNotBeNil)
				So(err.Error(),
					ShouldContainSubstring, "unable to parse link to URL")
				So(newurl, ShouldBeEmpty)
			})
		})
}

// Test FromHeadersOrDefaultGeneric
func Test_FromHeadersOrDefaultGeneric(t *testing.T) {
	Convey("Given HTTP headers with various forwarding configurations",
		t, func() {
			Convey("When creating a GenericBuilder from headers",
				func() {
					tests := []struct {
						description   string
						fwdHost       string
						fwdPathPrefix string
						expected      string
					}{
						{
							description:   "no forwarded headers are provided",
							fwdHost:       "",
							fwdPathPrefix: "",
							expected:      "http://localhost:8080",
						},
						{
							description: "forwarded host and path prefix " +
								"are provided",
							fwdHost:       "api.external.host",
							fwdPathPrefix: "prefix",
							expected: "https://api.external.host/" +
								"prefix",
						},
						{
							description: "only forwarded host with api. " +
								"prefix is provided",
							fwdHost:       "api.external.host",
							fwdPathPrefix: "",
							expected:      "https://api.external.host",
						},
						{
							description: "only forwarded path prefix is " +
								"provided",
							fwdHost:       "",
							fwdPathPrefix: "prefix",
							expected: "http://localhost:8080/" +
								"prefix",
						},
						{
							description: "internal forwarded host " +
								"without api. prefix is provided",
							fwdHost:       "internalhost",
							fwdPathPrefix: "",
							expected:      "http://localhost:8080",
						},
						{
							description: "internal forwarded host " +
								"without api. prefix and path prefix " +
								"are provided",
							fwdHost:       "internalhost",
							fwdPathPrefix: "prefix",
							expected: "http://localhost:8080/" +
								"prefix",
						},
					}

					for _, tt := range tests {
						Convey(tt.description, func() {
							// Given
							h := http.Header{}
							if tt.fwdHost != "" {
								h.Add("X-Forwarded-Host",
									tt.fwdHost)
							}
							if tt.fwdPathPrefix != "" {
								h.Add("X-Forwarded-Path-Prefix",
									tt.fwdPathPrefix)
							}

							// When
							builder :=
								FromHeadersOrDefaultGeneric(
									&h,
									defaultInternalURL,
									NoOpTransformer,
								)

							// Then
							So(builder, ShouldNotBeNil)
							So(builder.BaseURL.String(),
								ShouldEqual, tt.expected)
						})
					}
				})
		})
}

// Test RemovePrefixFromPath
func Test_RemovePrefixFromPath(t *testing.T) {
	Convey("Given a path and a prefix to remove", t, func() {
		Convey("When removing the prefix from the path", func() {
			tests := []struct {
				description string
				path        string
				prefix      string
				expected    string
			}{
				{
					description: "an empty path is provided",
					path:        "",
					prefix:      "/prefix",
					expected:    "",
				},
				{
					description: "a path with only the prefix is " +
						"provided",
					path:     "/prefix",
					prefix:   "/prefix",
					expected: "",
				},
				{
					description: "a path with the prefix and additional " +
						"segments is provided",
					path:     "/prefix/some/path",
					prefix:   "/prefix",
					expected: "/some/path",
				},
				{
					description: "a path with multiple prefix instances " +
						"is provided",
					path: "/prefix/prefix/prefix/" +
						"some/path",
					prefix:   "/prefix",
					expected: "/some/path",
				},
				{
					description: "a path without the prefix is " +
						"provided",
					path:     "/some/path",
					prefix:   "/prefix",
					expected: "/some/path",
				},
				{
					description: "a path without leading slash is " +
						"provided",
					path:     "some/path",
					prefix:   "/prefix",
					expected: "some/path",
				},
			}

			for _, tt := range tests {
				Convey(tt.description, func() {
					// When
					result := RemovePrefixFromPath(
						tt.path,
						tt.prefix,
					)

					// Then
					So(result, ShouldEqual, tt.expected)
				})
			}
		})
	})
}

// Test CustomTransformer
func Test_CustomTransformer(t *testing.T) {
	Convey("Given a custom path transformation function", t, func() {
		Convey("When building a link with custom transformation",
			func() {
				// Given
				baseURL, _ := url.Parse("https://example.com")

				customLogic := func(path string) string {
					// Example: prepend /api if not already there
					if !strings.HasPrefix(path, "/api") {
						return "/api" + path
					}
					return path
				}

				builder := NewGenericBuilder(
					baseURL,
					CustomTransformer(customLogic),
				)

				// When
				result, err := builder.BuildLink(
					"http://old.com/users",
				)

				// Then
				So(err, ShouldBeNil)
				So(result, ShouldEqual,
					"https://example.com/api/users")
			})
	})
}
