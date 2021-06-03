package api

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
)

type Sign struct {
	CallbackCheck SignCallbackCheck
}

type SignCallbackCheck func(r ...func(*SignCallbackCheckRequest)) bool

type SignCallbackCheckRequest struct {
	storeId  string
	passport string

	cb Callback

	ctx context.Context
}

func newSignCallbackCheckFunc(cfg Config) SignCallbackCheck {

	return func(o ...func(*SignCallbackCheckRequest)) bool {

		var (
			localSign  string
			r          SignCallbackCheckRequest
		)

		r = SignCallbackCheckRequest{
			storeId:  cfg.storeId,
			passport: cfg.passport,
		}

		for _, f := range o {
			f(&r)
		}

		localSign = makeSignature(
			[]byte(r.passport),
			[]byte(r.cb.StoreId),
			[]byte(r.cb.OrderId),
			[]byte(r.cb.PaymentState),
			[]byte(r.cb.Message),
			[]byte(r.passport),
		)

		return r.cb.StoreId == r.storeId && r.cb.Signature == localSign
	}
}

func (r SignCallbackCheck) WithContext(ctx context.Context) func(request *SignCallbackCheckRequest) {

	return func(r *SignCallbackCheckRequest) {
		r.ctx = ctx
	}
}

func (r SignCallbackCheck) WithParams(cb Callback) func(*SignCallbackCheckRequest) {

	return func(r *SignCallbackCheckRequest) {
		r.cb = cb
	}
}

func makeSignature(data ...[]byte) string {
	hasher := sha1.New()

	for _, d := range data {
		hasher.Write(d)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
