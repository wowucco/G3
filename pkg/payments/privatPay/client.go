package privatPay

import (
	"github.com/wowucco/G3/pkg/payments/privatPay/api"
	"github.com/wowucco/G3/pkg/payments/privatPay/transport"
	"net/http"
	"net/url"
)

type Config struct {
	StoreId     string
	Password    string
	Min         int
	Max         int
	ResponseUrl string
	RedirectUrl string

	Transport http.RoundTripper
	BaseUrl   *url.URL
}

func NewClient(cfg Config) *Client {

	tcfg := transport.Config{
		Transport: cfg.Transport,
		Url:       cfg.BaseUrl,
	}

	acfg := api.NewConfig(cfg.StoreId, cfg.Password, cfg.ResponseUrl, cfg.RedirectUrl)

	return &Client{
		min: cfg.Min,
		max: cfg.Max,
		API: api.New(transport.New(tcfg), acfg),
	}
}

type Client struct {
	min      int
	max      int

	*api.API
}
