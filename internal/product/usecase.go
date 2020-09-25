package product

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
)

type CreateProductForm struct {}

type UpdateProductFrom struct {}

type UseCase interface {
	Get(ctx context.Context, id int) (*entity.Product, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, offset, limit int) ([]*entity.Product, error)
	Create(ctx context.Context, form CreateProductForm) (*entity.Product, error)
	Update(ctx context.Context, id int, form UpdateProductFrom) (*entity.Product, error)
	Delete(ctx context.Context, id int) (*entity.Product, error)
}
