package strategy

import (
	"context"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/entity"
)

const providerNone = "none"

func NewDefaultStrategy(r checkout.IPaymentRepository) *DefaultStrategy {

	return &DefaultStrategy{r}
}

type DefaultStrategy struct {
	repository checkout.IPaymentRepository
}

func (s *DefaultStrategy) Init(ctx context.Context, order *entity.Order, payment *entity.Payment) (IInitPaymentStrategyResponse, error) {

	return NewIniPaymentStrategyResponse(entity.PaymentInitActionNone, "", providerNone), nil
}