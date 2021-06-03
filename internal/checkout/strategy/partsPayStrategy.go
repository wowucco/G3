package strategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/pkg/payments/privatPay"
	"github.com/wowucco/G3/pkg/payments/privatPay/api"
	"io/ioutil"
)

const providerPrivatPartsPay = "privat_parts_pay"
const providerPrivatPartsPayCtxKey = "privat_parts_pay_ctx_key"

func NewPartsPayStrategy(r checkout.IPaymentRepository, c *privatPay.Client) *PartsPayStrategy {

	return &PartsPayStrategy{r, c}
}

type PartsPayStrategy struct {
	repository checkout.IPaymentRepository
	provider   *privatPay.Client
}

func (s *PartsPayStrategy) Init(ctx context.Context, order *entity.Order, payment *entity.Payment) (IInitPaymentStrategyResponse, error) {

	tPrice := order.GetPrice()

	products := make([]api.Product, len(order.GetItems()))

	for k, p := range order.GetItems() {
		pPrice := p.GetPrice()
		products[k] = api.NewProduct(p.GetProduct().Name, p.GetQuantity(), pPrice.CentToFloatValue())
	}

	res, err := s.provider.Pay.Hold(
		s.provider.Pay.Hold.WithParams(payment.GetTransactionId(), tPrice.CentToFloatValue(), order.GetPayment().GetExtra().GetPartsPay(), products, api.MerchantTypePP),
		s.provider.Pay.Hold.WithContext(ctx),
	)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[privat pay init][hold payment][%v]", err))
	}

	defer res.Body.Close()

	var (
		result api.InitPaymentResponse
	)

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.New(fmt.Sprintf("[privat pay init][decode hold response][%v]", err))
	}

	if result.State != api.StateSuccess {
		return nil, errors.New(fmt.Sprintf("[privat pay init][response state is not success][%s][%s]", result.State, result.Message))
	}

	return NewIniPaymentStrategyResponse(entity.PaymentInitActionRedirect, api.PaymentRedirectUrl(result.Token), providerPrivatPartsPay), nil
}

func (s *PartsPayStrategy) Accept(ctx context.Context, order *entity.Order, payment *entity.Payment) (IAcceptHoldenPaymentStrategyResponse, error) {

	res, err := s.provider.Pay.Accept(s.provider.Pay.Accept.WithContext(ctx), s.provider.Pay.Accept.WithParams(payment.GetTransactionId()))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[privat pay accept][response][%v]", err))
	}

	defer res.Body.Close()

	var (
		status int
		desc   string
		result api.AcceptHoldenPaymentResponse
	)

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.New(fmt.Sprintf("[privat pay accept][decode response][%v]", err))
	}

	switch result.State {
	case api.StateFail:
		status = entity.PaymentStatusFailed
		desc = fmt.Sprintf("err_code: none,err_desc: %v", result.Message)
	case api.StateSuccess:
		status = entity.PaymentStatusConfirmed
		desc = "Holden payment was accept"
	default:
		return nil, errors.New(fmt.Sprintf("[privat pay accept][unhandled response status][%v]", result.State))
	}

	return NewAcceptHoldenPaymentStrategyResponse(status, desc, result.Stack()), nil
}
func (s *PartsPayStrategy) IsValidSignature(ctx *gin.Context) bool {

	cb, err := retrievePaymentCallback(ctx)

	if err != nil {
		return false
	}

	return s.provider.Sign.CallbackCheck(s.provider.Sign.CallbackCheck.WithParams(cb), s.provider.Sign.CallbackCheck.WithContext(ctx))
}
func (s *PartsPayStrategy) GetTransactionId(ctx *gin.Context) string {

	cb, err := retrievePaymentCallback(ctx)

	if err != nil {
		return ""
	}

	return cb.OrderId
}
func (s *PartsPayStrategy) ProcessingCallback(ctx *gin.Context) (IProcessingCallbackPaymentStrategyResponse, error) {

	var (
		status int
		desc   string
		skip   bool
	)
	cb, err := retrievePaymentCallback(ctx)

	if err != nil {
		return nil, err
	}

	switch cb.PaymentState {
	case api.StateSuccess:
		status = entity.PaymentStatusDone
		desc = "updated by callback"
	case api.StateFail:
		status = entity.PaymentStatusFailed
		desc = cb.Message
	case api.StateCanceled:
		status = entity.PaymentStatusCanceled
		desc = "updated by callback"
	case api.StateLocked:
		status = entity.PaymentStatusWaitingConfirmation
		desc = "updated by callback"
	case api.StateClientWait:
		fallthrough
	case api.StateOtpWaiting:
		fallthrough
	case api.StatePpCreation:
		fallthrough
	case api.StateCreated:
		skip = true
	default:
		return nil, errors.New(fmt.Sprintf("unknow privat part pay status '%s'", cb.PaymentState))
	}

	return NewProcessingCallbackPaymentStrategyResponse(status, desc, cb.Stack(), skip), nil
}

func retrievePaymentCallback(ctx *gin.Context) (api.Callback, error) {

	var cb api.Callback

	if d, exist := ctx.Get(providerPrivatPartsPayCtxKey); exist == true {

		cb = d.(api.Callback)
	} else {
		b, err := ioutil.ReadAll(ctx.Request.Body)

		if err != nil {
			return cb, errors.New(fmt.Sprintf("[failet to read request body][%v]", err))
		}

		if err := json.Unmarshal(b, &cb); err != nil {
			return cb, errors.New(fmt.Sprintf("[failed to unmarshal body][%v]", err))
		}

		ctx.Set(providerPrivatPartsPayCtxKey, cb)
	}

	return cb, nil
}
