package usecase

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/internal/product"
)

type ProductUseCaseMock struct {
	mock.Mock
}

func (m ProductUseCaseMock) Get(ctx context.Context, id int) (*entity.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m ProductUseCaseMock) Create(ctx context.Context, form product.CreateProductForm) (*entity.Product, error) {
	args := m.Called(form)

	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m ProductUseCaseMock) Update(ctx context.Context, id int, form product.UpdateProductFrom) (*entity.Product, error) {
	args := m.Called(form)

	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m ProductUseCaseMock) Delete(ctx context.Context, id int) (*entity.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Product), args.Error(1)
}