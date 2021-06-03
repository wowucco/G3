package api

import (
	"io"
	"net/http"
)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
}

type InitPaymentResponse struct {
	StoreId   string `json:"storeId"`
	State     string `json:"state"`
	OrderId   string `json:"orderId"`
	Token     string `json:"token"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (r InitPaymentResponse) Stack() map[string]interface{} {
	return map[string]interface{}{
		"state":     r.State,
		"storeId":   r.StoreId,
		"orderId":   r.OrderId,
		"token":     r.Token,
		"message":   r.Message,
		"signature": r.Signature,
	}
}

type AcceptHoldenPaymentResponse struct {
	StoreId   string `json:"storeId"`
	State     string `json:"state"`
	OrderId   string `json:"orderId"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (r AcceptHoldenPaymentResponse) Stack() map[string]interface{} {
	return map[string]interface{}{
		"state":     r.State,
		"orderId":   r.OrderId,
		"message":   r.Message,
		"storeId":   r.StoreId,
		"signature": r.Signature,
	}
}