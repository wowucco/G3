package delivery

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
)

// default cities list
	// get default ids from config etc
	// get list from es by ids
// search city
	// search city by text
// delivery info by city

type DeliveryReadRepository interface {
	// mix
	GetDeliveryInfoByCityId(ctx context.Context, id string) ([]*entity.DeliveryInfo, error)

	// Psql
	// GetDeliveryMethodsByCity(ctx context.Context, city entity.City) ([]*entity.DeliveryMethod, error)
	// GetPaymentMethodsByDeliveryMethodId(ctx context.Context, id uint) ([]*entity.PaymentMethod, error)

	// ES
	// GetDefaultCities(ctx context.Context) ([]*entity.City, error)
	GetCityById(ctx context.Context, id string) (*entity.City, error)
	SearchCity(ctx context.Context, text string) ([]*entity.City, error)
}
