package privatPay

import "github.com/wowucco/G3/pkg/payments/privatPay/api"

func NewClient(storeId, password string, min, max int) *Client {

	return &Client{
		storeId: storeId,
		password: password,
		min: min,
		max: max,
		API: api.New(),
	}
}

type Client struct {
	storeId string
	password string
	min int
	max int

	*api.API
}
