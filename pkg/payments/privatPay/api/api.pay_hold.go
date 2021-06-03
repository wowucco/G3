package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const holdUri = "/ipp/v2/payment/hold"

type PaymentHold func(r ...func(*PaymentHoldRequest)) (*Response, error)

type PaymentHoldRequest struct {
	storeId     string
	passport    string
	responseUrl string
	redirectUrl string

	orderId      string
	amount       float64
	partsCount   int
	merchantType string
	scheme       int
	products     []Product
	recipientId  string

	signature string

	ctx context.Context
}

func newPaymentHoldFunc(t Transport, cfg Config) PaymentHold {

	return func(o ...func(*PaymentHoldRequest)) (*Response, error) {
		var r = PaymentHoldRequest{
			storeId:     cfg.storeId,
			passport:    cfg.passport,
			responseUrl: cfg.responseUrl,
			redirectUrl: cfg.redirectUrl,
		}

		for _, f := range o {
			f(&r)
		}

		return r.Do(t)
	}
}

func (r PaymentHold) WithContext(ctx context.Context) func(request *PaymentHoldRequest) {

	return func(r *PaymentHoldRequest) {
		r.ctx = ctx
	}
}

func (r PaymentHold) WithParams(orderId string, amount float64, partsCount int, products []Product, merchantType string) func(*PaymentHoldRequest) {

	return func(r *PaymentHoldRequest) {
		r.orderId = orderId
		r.amount = amount
		r.partsCount = partsCount
		r.products = products
		r.merchantType = merchantType
	}
}

func (r PaymentHold) WithSchema(s int) func(*PaymentHoldRequest) {

	return func(r *PaymentHoldRequest) {
		r.scheme = s
	}
}

func (r PaymentHold) WithRecipientId(rid string) func(*PaymentHoldRequest) {

	return func(r *PaymentHoldRequest) {
		r.recipientId = rid
	}
}

func (r PaymentHoldRequest) Do(t Transport) (*Response, error) {
	var (
		buf            bytes.Buffer
		method         string
		params         map[string]interface{}
		products       []interface{}
		productsString string
	)

	method = "POST"

	for _, p := range r.products {
		products = append(products, map[string]interface{}{
			"name":  p.name,
			"count": p.count,
			"price": p.price,
		})

		productsString += fmt.Sprintf("%s%d%d", p.name, p.count, int(p.price) * 100)
	}

	params = map[string]interface{}{
		"storeId":      r.storeId,
		"orderId":      r.orderId,
		"amount":       r.amount,
		"partsCount":   r.partsCount,
		"merchantType": r.merchantType,
		"products":     products,
	}

	if r.scheme != 0 {
		params["scheme"] = r.scheme
	}

	if r.recipientId != "" {
		params["recipientId"] = r.recipientId
	}

	if r.responseUrl != "" {
		params["responseUrl"] = r.responseUrl
	}

	if r.redirectUrl != "" {
		params["redirectUrl"] = r.redirectUrl
	}

	params["signature"] = makeSignature(
		[]byte(r.passport),
		[]byte(r.storeId),
		[]byte(r.orderId),
		[]byte(strconv.Itoa(int(r.amount) * 100)),
		[]byte(strconv.Itoa(r.partsCount)),
		[]byte(r.merchantType),
		[]byte(r.responseUrl),
		[]byte(r.redirectUrl),
		[]byte(productsString),
		[]byte(r.passport),
	)

	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, errors.New(fmt.Sprintf("ppp failed encode body for hold request %v", err))
	}

	req, _ := newRequest(method, holdUri, &buf)

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
