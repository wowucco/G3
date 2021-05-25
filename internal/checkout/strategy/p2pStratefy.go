package strategy

import (
	"context"
	"errors"
	"fmt"
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
	default:
		return nil, errors.New(fmt.Sprintf("[unhandled liqpay status][%v]", r["status"]))
	}

	return NewAcceptHoldenPaymentStrategyResponse(status, desc, r), nil
}

func (s *P2PStrategy) IsValidSignature(data map[string]interface{}) bool {

	return s.provider.ValidateSign(data)
}
func (s *P2PStrategy) GetTransactionId(data map[string]interface{}) string {
	return s.provider.GetTransactionId(data)
}
func (s *P2PStrategy) ProcessingCallback(data map[string]interface{}) (IProcessingCallbackPaymentStrategyResponse, error) {
	
	cb, err := s.provider.Processing(data)

	if err != nil {
		return nil, err
	}

	status, err := mapLiqpayStatuses(cb.Status)

	if err != nil {
		return nil, err
	}

	return NewProcessingCallbackPaymentStrategyResponse(status, cb.Desc, cb.Stack), nil
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
