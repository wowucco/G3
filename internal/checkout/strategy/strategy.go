package strategy

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/entity"
)

type IInitPaymentStrategyResponse interface {
	GetAction() string
	GetResource() string
	GetProviderName() string
}

type IAcceptHoldenPaymentStrategyResponse interface {
	GetStatus() int
	GetDescription() string
	GetData() map[string]interface{}
}

type IProcessingCallbackPaymentStrategyResponse interface {
	GetStatus() int
	GetDescription() string
	GetData() map[string]interface{}
	Skip() bool
}

type IInitPaymentStrategy interface {
	Init(ctx context.Context, order *entity.Order, payment *entity.Payment) (IInitPaymentStrategyResponse, error)
}

type IAcceptHoldenStrategy interface {
	Accept(ctx context.Context, order *entity.Order, payment *entity.Payment) (IAcceptHoldenPaymentStrategyResponse, error)
}

type IProviderCallbackStrategy interface {
	IsValidSignature(ctx *gin.Context) bool
	GetTransactionId(ctx *gin.Context) string
	ProcessingCallback(ctx *gin.Context) (IProcessingCallbackPaymentStrategyResponse, error)
}

func NewPaymentContext(p2pStrategy *P2PStrategy, partsPayStrategy *PartsPayStrategy, defaultStrategy *DefaultStrategy) *PaymentContext {

	return &PaymentContext{p2pStrategy, partsPayStrategy, defaultStrategy}
}

type PaymentContext struct {
	// todo LOGGER
	p2pStrategy      *P2PStrategy
	partsPayStrategy *PartsPayStrategy
	defaultStrategy  *DefaultStrategy
}

func (c *PaymentContext) GetInitPaymentStrategy(method *entity.PaymentMethod) IInitPaymentStrategy {

	switch method.GetSlug() {
	case entity.PaymentMethodP2P:
		return c.p2pStrategy
	case entity.PaymentMethodPartsPay:
		return c.partsPayStrategy
	default:
		return c.defaultStrategy
	}
}

func (c *PaymentContext) GetAcceptHoldenPaymentStrategy(provider string) (IAcceptHoldenStrategy, error) {
	switch provider {
	case providerLiqpay:
		return c.p2pStrategy, nil
	case providerPrivatPartsPay:
		return c.partsPayStrategy, nil
	default:
		return nil, errors.New(fmt.Sprintf("provider %s don't have accept interface", provider))
	}
}

func (c *PaymentContext) GetProviderCallbackPaymentStrategy(provider string) (IProviderCallbackStrategy, error) {
	switch provider {
	case providerLiqpay:
		return c.p2pStrategy, nil
	case providerPrivatPartsPay:
		return c.partsPayStrategy, nil
	default:
		return nil, errors.New(fmt.Sprintf("provider %s don't have accept interface", provider))
	}
}

func NewIniPaymentStrategyResponse(action, resource, provider string) IInitPaymentStrategyResponse {

	return &InitPaymentStrategyResponse{
		action:   action,
		resource: resource,
		provider: provider,
	}
}

type InitPaymentStrategyResponse struct {
	payment  *entity.Payment
	order    *entity.Order
	action   string
	resource string
	provider string
}

func (r *InitPaymentStrategyResponse) GetAction() string {
	return r.action
}
func (r *InitPaymentStrategyResponse) GetResource() string {
	return r.resource
}
func (r *InitPaymentStrategyResponse) GetProviderName() string {
	return r.provider
}

type AcceptHoldenPaymentStrategyResponse struct {
	status int
	desc   string
	stack  map[string]interface{}
}

func NewAcceptHoldenPaymentStrategyResponse(status int, desc string, stack map[string]interface{}) IAcceptHoldenPaymentStrategyResponse {

	return &AcceptHoldenPaymentStrategyResponse{status, desc, stack}
}
func (r *AcceptHoldenPaymentStrategyResponse) GetStatus() int {
	return r.status
}
func (r *AcceptHoldenPaymentStrategyResponse) GetDescription() string {
	return r.desc
}
func (r *AcceptHoldenPaymentStrategyResponse) GetData() map[string]interface{} {
	return r.stack
}

type ProcessingCallbackPaymentStrategyResponse struct {
	status int
	desc   string
	stack  map[string]interface{}
	skip   bool
}

func NewProcessingCallbackPaymentStrategyResponse(status int, desc string, stack map[string]interface{}, skip bool) IProcessingCallbackPaymentStrategyResponse {

	return &ProcessingCallbackPaymentStrategyResponse{status, desc, stack, skip}
}
func (r *ProcessingCallbackPaymentStrategyResponse) GetStatus() int {
	return r.status
}
func (r *ProcessingCallbackPaymentStrategyResponse) GetDescription() string {
	return r.desc
}
func (r *ProcessingCallbackPaymentStrategyResponse) GetData() map[string]interface{} {
	return r.stack
}
func (r *ProcessingCallbackPaymentStrategyResponse) Skip() bool {
	return r.skip
}
