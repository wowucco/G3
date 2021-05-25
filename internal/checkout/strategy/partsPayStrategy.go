package strategy

import (
	"context"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/entity"
)

const providerPrivatPartsPay = "privatPartsPay"

func NewPartsPayStrategy(r checkout.IPaymentRepository) *PartsPayStrategy {

	return &PartsPayStrategy{r}
}

type PartsPayStrategy struct {
	repository checkout.IPaymentRepository
}

func (s *PartsPayStrategy) Init(ctx context.Context, order *entity.Order, payment *entity.Payment) (IInitPaymentStrategyResponse, error) {

	return nil, nil
}

func (s *PartsPayStrategy) Accept(ctx context.Context, order *entity.Order, payment *entity.Payment) (IAcceptHoldenPaymentStrategyResponse, error) {

	return nil, nil
}

func (s *PartsPayStrategy) IsValidSignature(map[string]interface{}) bool {

	return false
}
func (s *PartsPayStrategy) GetTransactionId(map[string]interface{}) string {
	return ""
}
func (s *PartsPayStrategy) ProcessingCallback(map[string]interface{}) (IProcessingCallbackPaymentStrategyResponse, error) {

	return nil, nil
}