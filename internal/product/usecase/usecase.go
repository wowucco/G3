package usecase

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
	_product "github.com/wowucco/G3/internal/product"
)

type ProductUseCase struct {
	productRepo _product.Repository
}

func NewProductUseCase(productRepo _product.Repository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

func (p ProductUseCase) Get(ctx context.Context, id int) (*entity.Product, error) {
	product, err := p.productRepo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p ProductUseCase) Count(ctx context.Context) (int, error) {
	return p.productRepo.Count(ctx)
}

func (p ProductUseCase) Query(ctx context.Context, offset, limit int) ([]*entity.Product, error) {

	return p.productRepo.Query(ctx, offset, limit)
}

func (p ProductUseCase) Create(ctx context.Context, form _product.CreateProductForm) (*entity.Product, error) {

	product := &entity.Product{}

	id, err := p.productRepo.Create(ctx, product)

	if err != nil {
		return nil, err
	}

	product, err = p.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p ProductUseCase) Update(ctx context.Context, id int, form _product.UpdateProductFrom) (*entity.Product, error) {

	product, err := p.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if err = p.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (p ProductUseCase) Delete(ctx context.Context, id int) (*entity.Product, error) {

	product, err := p.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if err = p.productRepo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return product, nil
}