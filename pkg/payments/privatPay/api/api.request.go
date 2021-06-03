package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request interface {
	Do(t Transport) (*Response, error)
}

func newRequest(method, uri string, body io.Reader) (*http.Request, error) {

	r := http.Request{
		Method:     method,
		URL:        &url.URL{Path: uri},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}

	if body != nil {
		switch b := body.(type) {
		case *bytes.Buffer:
			r.Body = ioutil.NopCloser(body)
			r.ContentLength = int64(b.Len())
		case *bytes.Reader:
			r.Body = ioutil.NopCloser(body)
			r.ContentLength = int64(b.Len())
		case *strings.Reader:
			r.Body = ioutil.NopCloser(body)
			r.ContentLength = int64(b.Len())
		default:
			r.Body = ioutil.NopCloser(body)
		}
	}

	return &r, nil
}

type Callback struct {
	OrderId      string `json:"orderId"`
	StoreId      string `json:"storeId"`
	PaymentState string `json:"paymentState"`
	Message      string `json:"message"`
	Signature    string `json:"signature"`
}

func (r Callback) Stack() map[string]interface{} {
	return map[string]interface{}{
		"paymentState": r.PaymentState,
		"orderId":      r.OrderId,
		"storeId":      r.StoreId,
		"message":      r.Message,
		"signature":    r.Signature,
	}
}