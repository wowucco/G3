package repository

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
	"log"
)

func (d DeliveryReadRepository) GetDeliveryInfoByCityId(ctx context.Context, id string) ([]*entity.DeliveryInfo, error) {

	city, err := d.GetCityById(ctx, id)

	if err != nil {
		return nil, err
	}

	deliveryMethods, err := d.db.getDeliveryMethodsByCity(ctx, *city)

	if err != nil {
		return nil, err
	}

	deliveryInfo := make([]*entity.DeliveryInfo, len(deliveryMethods))

	for key, deliveryMethod := range deliveryMethods {

		paymentMethods, _ := d.db.getPaymentMethodsByDeliveryMethodId(ctx, deliveryMethod.ID)
		warehouses, _ := d.getWarehousesOfCityByDeliveryMethod(ctx, *city, deliveryMethod)

		deliveryInfo[key] = &entity.DeliveryInfo{
			DeliveryMethod: deliveryMethod,
			PaymentMethods: paymentMethods,
			Warehouses:     warehouses,
		}
	}

	return deliveryInfo, nil
}

func (d DeliveryReadRepository) GetCityById(ctx context.Context, id string) (*entity.City, error) {

	return d.es.getCityById(ctx, id)
}

func (d DeliveryReadRepository) SearchCity(ctx context.Context, text string) ([]*entity.City, error) {

	return d.es.searchCity(ctx, text)
}

func (d DeliveryReadRepository) GetDeliveryMethodBySlug(slug string) (*entity.DeliveryMethod, error) {

	return d.db.getDeliveryMethodBySlug(slug)
}

func (d DeliveryReadRepository) GetPaymentMethodBySlug(slug string) (*entity.PaymentMethod, error) {

	return d.db.getPaymentMethodBySlug(slug)
}

func (d DeliveryReadRepository) getWarehousesOfCityByDeliveryMethod(ctx context.Context, city entity.City, deliveryMethod entity.DeliveryMethod) ([]entity.Warehouse, error) {

	switch deliveryMethod.Slug {
	case entity.DeliveryMethodYourself:
		return d.db.getWarehousesForYourselfByCity(ctx, city)
	case entity.DeliveryMethodNovaposhta:
		return d.getWarehousesForNovaposhtaByCity(ctx, city)
	case entity.DeliveryMethodCourier:
		fallthrough
	default:
		return nil, nil
	}
}

func (d DeliveryReadRepository) getWarehousesForNovaposhtaByCity(ctx context.Context, city entity.City) ([]entity.Warehouse, error) {

	esr, err := d.es.getWarehousesForNovaposhtaByCity(ctx, city)

	if err != nil {
		return make([]entity.Warehouse, 0), err
	}
	
	if len(esr) > 0 {
		return toWarehouseEntities(esr), nil
	}

	npwh, err := d.np.getWarehousesForNovaposhtaByCity(ctx, city)

	if err != nil {
		return make([]entity.Warehouse, 0), err
	}

	go func() {
		if e := d.es.reindexCitiesWarehouses(city, npwh); e != nil {
			log.Fatalf("reindex warehouses error: %s", err)
		}
	}()

	return toWarehouseEntities(npwh), nil
}

func toWarehouseEntities(npwh []NPWarehouse) []entity.Warehouse {

	w := make([]entity.Warehouse, len(npwh))

	for k, v := range npwh {
		w[k] = entity.Warehouse{
			ID:      v.ID,
			Name:    v.NameRu,
			Address: v.ShortAddressRu,
			Phone:   v.Phone,
		}
	}

	return w
}
