package product

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
)

type Repository interface {
	Get(ctx context.Context, id int) (*entity.Product, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, offset, limit int) ([]*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) (int, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int) error
}