package checkout

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
)

type IOrderRepository interface {
	NextId() (int, error)
	Get(ctx context.Context, orderId int) (*entity.Order, error)
	Save(ctx context.Context, order *entity.Order) error
	Create(ctx context.Context, builder *CreateOrderBuilder) (*entity.Order, error)
}

type IPaymentRepository interface {
	NextId() (int, error)
	Get(ctx context.Context, transactionId string) (*entity.Payment, error)
	Save(ctx context.Context, p *entity.Payment) error
	Create(ctx context.Context, p *entity.Payment) error
}
