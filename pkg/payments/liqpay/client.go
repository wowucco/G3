package liqpay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	_liqpay "github.com/liqpay/go-sdk"
	"html/template"
	"log"
)

//private const LIQPAY_API_PATH = 'request';
//
//private const ACTION_HOLD = 'hold';
//private const ACTION_HOLD_COMPLETION = 'hold_completion';
//
//private const PAYMENT_STATUS_HOLD_WAIT = 'hold_wait';
//
//private const PAYMENT_STATUS_ERROR   = 'error';
//private const PAYMENT_STATUS_SUCCESS = 'success';
//private const PAYMENT_STATUS_FAILURE = 'failure';

const (
	path        = "request"
	version     = "3"
	currencyUah = "UAH"

	actionHold       = "hold"
	actionAcceptHold = "hold_completion"

	StatusHoldWait = "hold_wait"

	StatusError   = "error"
	StatusSuccess = "success"
	StatusFail    = "failure"
)

type CallbackResponse struct {
	Status string
	Desc   string
	Stack  map[string]interface{}
}

func NewClient(publicKey, privateKey, callbackUrl, returnUrl string) *Client {

	return &Client{
		api:         _liqpay.New(publicKey, privateKey, nil),
		callbackUrl: callbackUrl,
		returnUrl:   returnUrl,
		publicKey:   publicKey,
	}
}

type Client struct {
	api *_liqpay.Client

	callbackUrl string
	returnUrl   string

	publicKey string
}

func (c *Client) Hold(orderId, amount, description string) (string, error) {

	r := _liqpay.Request{
		"public_key":  c.publicKey,
		"order_id":    orderId,
		"amount":      amount,
		"description": description,
		"action":      actionHold,
		"version":     version,
		"currency":    currencyUah,
		"result_url":  c.returnUrl,
		"server_url":  c.callbackUrl,
	}

	log.Printf("[p2p init][liqpay hold response][%s][%v]", orderId, r)

	return c.RenderForm(r)
}

type formData struct {
	Data      string
	Signature string
}

// todo
// go get github.com/liqpay/go-sdk can't pull liqpay_form.html file
func (c Client) RenderForm(req _liqpay.Request) (string, error) {

	encodedJSON, err := req.Encode()
	if err != nil {
		return "", err
	}

	signature := c.api.Sign([]byte(encodedJSON))

	t, err := template.ParseFiles("./pkg/payments/liqpay/liqpay_form.html")

	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	if err := t.Execute(&buf, formData{
		Data:      encodedJSON,
		Signature: signature,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Client) AcceptHolden(orderId, amount string) (map[string]interface{}, error) {

	r := _liqpay.Request{
		"action":   actionAcceptHold,
		"version":  version,
		"order_id": orderId,
		"amount":   amount,
	}

	response, err := c.api.Send(path, r)

	if err != nil {
		log.Printf("[p2p accept holden][liqpay request error][%s][%v][%v]", orderId, r, response)
	} else {
		log.Printf("[p2p accept holden][liqpay request][%s][%v][%v]", orderId, r, response)
	}

	return response, err
}

func (c *Client) ValidateSign(data map[string]interface{}) bool {
	return c.api.Sign([]byte(data["data"].(string))) == data["signature"].(string)
}

func (c *Client) GetTransactionId(data map[string]interface{}) string {

	b := c.decodeData(data["data"].(string))

	return b["order_id"].(string)
}

func (c *Client) Processing(data map[string]interface{}) (*CallbackResponse, error) {

	b := c.decodeData(data["data"].(string))

	switch b["action"].(string) {
	case actionHold:
		return c.handleByHoldCallback(b)
	}

	return nil, errors.New(fmt.Sprintf("unhandled provider action %v", b["action"]))
}

func (c *Client) decodeData(data string) map[string]interface{} {
	d, _ := base64.StdEncoding.DecodeString(data)

	body := make(map[string]interface{})

	_ = json.Unmarshal(d, &body)

	return body
}

func (c *Client) handleByHoldCallback(body map[string]interface{}) (*CallbackResponse, error) {

	cb := &CallbackResponse{
		Stack: body,
	}

	switch body["status"].(string) {
	case StatusHoldWait:
		cb.Status = StatusHoldWait
		return cb, nil
	case StatusError:
		cb.Status = StatusError
		cb.Desc = fmt.Sprintf("code: %v | description: %v", body["err_code"], body["err_description"])
		return cb, nil
	case StatusSuccess:
		cb.Status = StatusSuccess
		return cb, nil
	default:
		return cb, errors.New(fmt.Sprintf("unhandled status '%v' for hold action", body["status"]))
	}
}
