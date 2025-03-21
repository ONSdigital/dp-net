package http

//go:generate moq -out mock_client.go . Clienter

import (
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	request "github.com/ONSdigital/dp-net/v3/request"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// go:generate moq -out mock_client.go -pkg http . Clienter

const (
	DefaultRequestTimeout = 10 * time.Second
)

// Client is an extension of the net/http client with ability to add
// timeouts, exponential backoff and context-based cancellation.
type Client struct {
	MaxRetries         int
	RetryTime          time.Duration
	PathsWithNoRetries map[string]bool
	HTTPClient         *http.Client
	TotalTimeout       time.Duration
}

// DefaultTransport is the default implementation of Transport and is
// used by DefaultClient.
var DefaultTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout: 5 * time.Second,
	MaxIdleConns:        10,
	IdleConnTimeout:     30 * time.Second,
}

// DefaultClient is a dp-net specific http client with sensible timeouts,
// exponential backoff, and a contextual dialer.
// NOTE: This is needed in dp-deployer's unit tests
var DefaultClient = &Client{
	MaxRetries: 3,
	RetryTime:  20 * time.Millisecond,

	HTTPClient: &http.Client{
		Timeout:   10 * time.Second,
		Transport: DefaultTransport,
	},
}

// Clienter provides an interface for methods on an HTTP Client.
type Clienter interface {
	// Sets the overall timeout (spans all retries)
	SetTotalTimeout(timeout time.Duration)
	// Sets HTTP request timeout (timeout 'per' try/request)
	SetTimeout(timeout time.Duration)
	SetMaxRetries(int)
	GetMaxRetries() int
	SetPathsWithNoRetries([]string)
	GetPathsWithNoRetries() []string

	Get(ctx context.Context, url string) (*http.Response, error)
	Head(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error)
	Put(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error)
	PostForm(ctx context.Context, uri string, data url.Values) (*http.Response, error)

	Do(ctx context.Context, req *http.Request) (*http.Response, error)
	RoundTrip(req *http.Request) (*http.Response, error)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewClient returns a copy of DefaultClient.
func NewClient() Clienter {
	return &Client{
		MaxRetries: 3,
		RetryTime:  20 * time.Millisecond,

		HTTPClient: &http.Client{
			Timeout:   DefaultRequestTimeout,
			Transport: DefaultTransport,
		},
	}
}

// NewClientWithTransport returns a copy of default client, plus changes to the transport layer
func NewClientWithTransport(transport http.RoundTripper) Clienter {
	client := &Client{
		MaxRetries: 3,
		RetryTime:  20 * time.Millisecond,

		HTTPClient: &http.Client{
			Timeout:   DefaultRequestTimeout,
			Transport: transport,
		},
	}

	return client
}

// ClientWithTimeout creates a client and sets per try/request timeout.
func ClientWithTimeout(c Clienter, timeout time.Duration) Clienter {
	if c == nil {
		c = NewClient()
	}
	c.SetTimeout(timeout)
	return c
}

// ClientWithTotalTimeout creates a client with overall timeout (spans all retries)
func ClientWithTotalTimeout(c Clienter, timeout time.Duration) Clienter {
	if c == nil {
		c = NewClient()
	}
	c.SetTotalTimeout(timeout)
	return c
}

// ClientWithTimeouts creates a client with both (per-request + total) timeout values
func ClientWithTimeouts(c Clienter, perRequest time.Duration, total time.Duration) Clienter {
	if c == nil {
		c = NewClient()
	}
	c.SetTimeout(perRequest)
	c.SetTotalTimeout(total)
	return c
}

// Clienter roundtripper calls the httpclient roundtripper.
func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Transport.RoundTrip(req)
}

// ClientWithListOfNonRetriablePaths facilitates creating a client and setting a
// list of paths that should not be retried on failure.
func ClientWithListOfNonRetriablePaths(c Clienter, paths []string) Clienter {
	if c == nil {
		c = NewClient()
	}
	c.SetPathsWithNoRetries(paths)
	return c
}

// SetTotalTimeout sets the overall timeout (spans all retries)
func (c *Client) SetTotalTimeout(timeout time.Duration) {
	c.TotalTimeout = timeout
}

// SetTimeout sets HTTP request timeout (timeout 'per' try/request)
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

// GetMaxRetries gets the HTTP request maximum number of retries.
func (c *Client) GetMaxRetries() int {
	return c.MaxRetries
}

// SetMaxRetries sets HTTP request maximum number of retries.
func (c *Client) SetMaxRetries(maxRetries int) {
	c.MaxRetries = maxRetries
}

// GetPathsWithNoRetries gets a list of paths that will HTTP request will not retry on error.
func (c *Client) GetPathsWithNoRetries() (paths []string) {
	for path, _ := range c.PathsWithNoRetries {
		paths = append(paths, path)
	}
	return paths
}

// SetPathsWithNoRetries sets a list of paths that will HTTP request will not retry on error.
func (c *Client) SetPathsWithNoRetries(paths []string) {
	mapPath := make(map[string]bool)
	for _, path := range paths {
		mapPath[path] = true
	}
	c.PathsWithNoRetries = mapPath
}

// Do calls ctxhttp.Do with the addition of retries with exponential backoff
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {

	// TODO: Remove this once user token (Florence token) is propegated throughout apps
	// Used for audit purposes
	if request.IsUserPresent(ctx) {
		// only add this header if not already set
		if len(req.Header.Get(request.UserHeaderKey)) == 0 {
			request.AddUserHeader(req, request.User(ctx))
		}
	}

	// get any existing correlation-id (might be "id1,id2"), append a new one, add to headers
	upstreamCorrelationIDs := request.GetRequestId(ctx)
	addedIDLen := 20
	if upstreamCorrelationIDs != "" {
		// get length of (first of) IDs (e.g. "id1" is 3), new ID will be half that size
		addedIDLen = len(upstreamCorrelationIDs) / 2
		if commaPosition := strings.Index(upstreamCorrelationIDs, ","); commaPosition > 1 {
			addedIDLen = commaPosition / 2
		}
		upstreamCorrelationIDs += ","
	}
	request.AddRequestIdHeader(req, upstreamCorrelationIDs+request.NewRequestID(addedIDLen))

	doer := func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		if req.ContentLength > 0 {
			var err error
			req.Body, err = req.GetBody()
			if err != nil {
				return nil, err
			}
		}
		return ctxhttp.Do(ctx, client, req)
	}

	path := req.URL.Path

	// A global timeout value is defined
	if c.TotalTimeout > 0 {
		// if there is a timeout set already, honor that
		// will allow user to set a local/per-request override
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.TotalTimeout)
			defer cancel()
		}
	}
	resp, err := doer(ctx, c.HTTPClient, req)
	if !c.PathsWithNoRetries[path] && c.GetMaxRetries() > 0 && wantRetry(err, resp) {
		return c.backoff(ctx, doer, c.HTTPClient, req)
	}

	return resp, err
}

func wantRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}
	if resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusConflict {
		return true
	}
	return false
}

// Get calls Do with a GET.
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req)
}

// Head calls Do with a HEAD.
func (c *Client) Head(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req)
}

// Post calls Do with a POST and the appropriate content-type and body.
func (c *Client) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.Do(ctx, req)
}

// Put calls Do with a PUT and the appropriate content-type and body.
func (c *Client) Put(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.Do(ctx, req)
}

// PostForm calls Post with the appropriate form content-type.
func (c *Client) PostForm(ctx context.Context, uri string, data url.Values) (*http.Response, error) {
	return c.Post(ctx, uri, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

type Doer = func(context.Context, *http.Client, *http.Request) (*http.Response, error)

func (c *Client) backoff(
	ctx context.Context,
	doer Doer,
	client *http.Client,
	req *http.Request,
) (resp *http.Response, err error) {

	for retries := 1; retries <= c.GetMaxRetries(); retries++ {
		pingChan := make(chan struct{}, 0)
		go func() {
			time.Sleep(getSleepTime(retries, c.RetryTime))
			close(pingChan)
		}()
		// check for first of: context cancellation or sleep ends
		select {
		case <-pingChan:
		case <-ctx.Done():
			err = ctx.Err()
			return
		}

		resp, err = doer(ctx, client, req)
		// prioritise any context cancellation
		if ctx.Err() != nil {
			err = ctx.Err()
			return
		}
		if !wantRetry(err, resp) {
			return
		}
	}
	return
}

// getSleepTime will return a sleep time based on the attempt and initial retry time.
// It uses the algorithm 2^n where n is the attempt number (double the previous) and
// a randomization factor of between 0-5ms so that the server isn't being hit constantly
// at the same time by many clients.
func getSleepTime(attempt int, retryTime time.Duration) time.Duration {
	n := (math.Pow(2, float64(attempt)))
	rnd := time.Duration(rand.Intn(4)+1) * time.Millisecond
	return (time.Duration(n) * retryTime) - rnd
}
