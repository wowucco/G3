package checkout

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/entity"
)

type IInitPaymentResponse interface {
	GetAction() string
	GetResource() string
	GetPaymentTransactionID() string
	GetPaymentMethod() string
	GetOrderId() int
	GetDoNotCall() bool
}

type IProviderCallbackPaymentResponse interface {
	GetStatusCode() int
	GetBody() map[string]interface{}
}

type IOrderUseCase interface {
	Create(ctx context.Context, form CreateOrderForm) (*entity.Order, error)
	InitPayment(ctx context.Context, form InitPaymentForm) (IInitPaymentResponse, error)
	AcceptHoldenPayment(ctx context.Context, form IAcceptHoldenPaymentForm) error
	ProviderCallback(ctx *gin.Context, form IProviderCallbackPaymentForm) (IProviderCallbackPaymentResponse, error)
}
