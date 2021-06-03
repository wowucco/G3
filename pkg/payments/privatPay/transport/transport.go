package transport

import (
	"net/http"
	"net/url"
	"strings"
)

const baseUrl = "https://payparts2.privatbank.ua"

type Config struct {
	Transport http.RoundTripper
	Url       *url.URL
}

func New(cfg Config) *Client {

	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}

	if cfg.Url == nil {
		cfg.Url, _ = url.Parse(baseUrl)
	}

	return &Client{
		transport: cfg.Transport,
		url:       cfg.Url,
	}
}

type Client struct {
	transport http.RoundTripper
	url       *url.URL
}

func (c *Client) Perform(req *http.Request) (*http.Response, error) {

	baseUrl := c.url

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	req.URL.Scheme = baseUrl.Scheme
	req.URL.Host = baseUrl.Host

	if baseUrl.Path != "" {
		var b strings.Builder
		b.Grow(len(baseUrl.Path) + len(req.URL.Path))
		b.WriteString(baseUrl.Path)
		b.WriteString(req.URL.Path)
		req.URL.Path = b.String()
	}

	return c.transport.RoundTrip(req)
}
