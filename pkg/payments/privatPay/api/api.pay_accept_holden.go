package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const acceptHoldenUri = "/ipp/v2/payment/confirm"

type PaymentAcceptHolden func(r ...func(*PaymentAcceptHoldenRequest)) (*Response, error)

type PaymentAcceptHoldenRequest struct {
	storeId  string
	passport string

	orderId string

	signature string

	ctx context.Context
}

func newAcceptHoldenFunc(t Transport, cfg Config) PaymentAcceptHolden {

	return func(o ...func(*PaymentAcceptHoldenRequest)) (*Response, error) {
		var r = PaymentAcceptHoldenRequest{
			storeId:  cfg.storeId,
			passport: cfg.passport,
		}

		for _, f := range o {
			f(&r)
		}

		return r.Do(t)
	}
}

func (r PaymentAcceptHolden) WithContext(ctx context.Context) func(request *PaymentAcceptHoldenRequest) {

	return func(r *PaymentAcceptHoldenRequest) {
		r.ctx = ctx
	}
}

func (r PaymentAcceptHolden) WithParams(orderId string) func(*PaymentAcceptHoldenRequest) {

	return func(r *PaymentAcceptHoldenRequest) {
		r.orderId = orderId
	}
}

func (r PaymentAcceptHoldenRequest) Do(t Transport) (*Response, error) {
	var (
		buf    bytes.Buffer
		method string
		params map[string]interface{}
	)

	method = "POST"

	params = map[string]interface{}{
		"storeId": r.storeId,
		"orderId": r.orderId,
	}

	params["signature"] = makeSignature(
		[]byte(r.passport),
		[]byte(r.storeId),
		[]byte(r.orderId),
		[]byte(r.passport),
	)

	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, errors.New(fmt.Sprintf("ppp failed encode body for request %v", err))
	}

	req, _ := newRequest(method, acceptHoldenUri, &buf)

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	res, err := t.Perform(req)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("ppp failed hold request %v", err))
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}
