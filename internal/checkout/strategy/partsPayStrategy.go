package strategy

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (s *PartsPayStrategy) RetrieveData(ctx *gin.Context) {

}

func (s *PartsPayStrategy) IsValidSignature(ctx *gin.Context) bool {

	return false
}
func (s *PartsPayStrategy) GetTransactionId(ctx *gin.Context) string {
	return ""
}
func (s *PartsPayStrategy) ProcessingCallback(ctx *gin.Context) (IProcessingCallbackPaymentStrategyResponse, error) {

	return nil, nil
}