package thunderapi

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/rs/zerolog"
)

const (
	defaultBaseURL = "https://api.thunderatz.org/"
	userAgent      = "go-thunderapi"
)

// A Client manages communication with ThundeRatz API
type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	logger zerolog.Logger

	BaseURL   *url.URL
	UserAgent string
}

// NewClient returns a new ThunderAPI client. If a nil httpClient is
// provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	return c
}
