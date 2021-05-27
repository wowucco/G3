package strategy

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/pkg/payments/liqpay"
)

const providerLiqpay = "liqpay"

func NewP2PStrategy(r checkout.IPaymentRepository, p *liqpay.Client) *P2PStrategy {

	return &P2PStrategy{r, p}
}

type P2PStrategy struct {
	repository checkout.IPaymentRepository

	provider *liqpay.Client
}

func (s *P2PStrategy) Init(ctx context.Context, order *entity.Order, payment *entity.Payment) (IInitPaymentStrategyResponse, error) {

	form, err := s.provider.Hold(payment.GetTransactionId(), payment.GetPrice().CentToCurrency(), payment.GetDescription())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[p2p init][hold payment][%v]", err))
	}

	return NewIniPaymentStrategyResponse(entity.PaymentInitActionForm, form, providerLiqpay), nil
}

func (s *P2PStrategy) Accept(ctx context.Context, order *entity.Order, payment *entity.Payment) (IAcceptHoldenPaymentStrategyResponse, error) {

	r, e := s.provider.AcceptHolden(payment.GetTransactionId(), payment.GetPrice().CentToCurrency())

	if e != nil {
		return nil, errors.New(fmt.Sprintf("[p2p accept error][%v]", e))
	}

	var (
		status int
		desc string
	)

	switch r["status"].(string) {
	case liqpay.StatusError:
		status = entity.PaymentStatusFailed
		desc = fmt.Sprintf("err_code: %v, err_desc: %v", r["err_code"], r["err_decription"])
	case liqpay.StatusFail:
		status = entity.PaymentStatusFailed
		desc = fmt.Sprintf("err_code: %v, err_desc: %v", r["err_code"], r["err_decription"])
	case liqpay.StatusSuccess:
		status = entity.PaymentStatusDone
		desc = "Holden payment was accept"
	case liqpay.StatusReversed:
		status = entity.PaymentStatusRefund
		desc = "Holden payment was reversed"
	default:
		return nil, errors.New(fmt.Sprintf("[unhandled liqpay status][%v]", r["status"]))
	}

	return NewAcceptHoldenPaymentStrategyResponse(status, desc, r), nil
}

func (s *P2PStrategy) IsValidSignature(ctx *gin.Context) bool {
	cb := make(map[string]interface{})
	cb["data"] = ctx.PostForm("data")
	cb["signature"] = ctx.PostForm("signature")
	return s.provider.ValidateSign(cb)
}
func (s *P2PStrategy) GetTransactionId(ctx *gin.Context) string {
	cb := make(map[string]interface{})
	cb["data"] = ctx.PostForm("data")
	cb["signature"] = ctx.PostForm("signature")
	return s.provider.GetTransactionId(cb)
}
func (s *P2PStrategy) ProcessingCallback(ctx *gin.Context) (IProcessingCallbackPaymentStrategyResponse, error) {
	cb := make(map[string]interface{})
	cb["data"] = ctx.PostForm("data")
	cb["signature"] = ctx.PostForm("signature")

	res, err := s.provider.Processing(cb)

	if err != nil {
		return nil, err
	}

	status, err := mapLiqpayStatuses(res.Status)

	if err != nil {
		return nil, err
	}

	return NewProcessingCallbackPaymentStrategyResponse(status, res.Desc, res.Stack), nil
}

func mapLiqpayStatuses(s string) (int, error) {
	switch s {
	case liqpay.StatusHoldWait:
		return entity.PaymentStatusWaitingConfirmation, nil
	case liqpay.StatusSuccess:
		return entity.PaymentStatusDone, nil
	case liqpay.StatusError:
		return entity.PaymentStatusFailed, nil
	case liqpay.StatusFail:
		return entity.PaymentStatusFailed, nil
	default:
		return 0, errors.New(fmt.Sprintf("unknow liqpay status '%s'", s))
	}
}
