package graph

import (
	"github.com/wowucco/G3/internal/delivery"
	"github.com/wowucco/G3/internal/menu"
	"github.com/wowucco/G3/internal/product"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	useCase product.UseCase
	productRead product.ReadRepository
	menuRead menu.ReadRepository
	deliveryRead delivery.DeliveryReadRepository
}
