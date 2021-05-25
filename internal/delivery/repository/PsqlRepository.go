package repository

import (
	"context"
	"encoding/json"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/entity"
)

const tableNameDeliveryMethods = "shop_delivery_method"
const tableNamePaymentMethods = "shop_payment_method"
const tableNameDeliveryAssignCity = "shop_delivery_assignment_city"
const tableNameDeliveryAssignPayment = "shop_delivery_assignment_payment"

func (d PsqlDeliveryReadRepository) getDeliveryMethodsByCity(ctx context.Context, city entity.City) ([]entity.DeliveryMethod, error) {

	var rows []DeliveryMethod

	err := d.db.Select("d.*").
		From(tableWithAlias(tableNameDeliveryMethods, "d")).
		LeftJoin(tableWithAlias(tableNameDeliveryAssignCity, "a"), dbx.NewExp("d.id = a.delivery_method_id")).
		Where(dbx.NewExp("d.status={:status}", dbx.Params{"status": true})).
		AndWhere(dbx.Or(dbx.NewExp("d.alpha={:alpha}", dbx.Params{"alpha": true}), dbx.NewExp("a.city_token={:token}", dbx.Params{"token": city.ID}))).
		OrderBy("d.width desc").
		All(&rows)

	if err != nil {
		return nil, err
	}

	deliveryMethods := make([]entity.DeliveryMethod, len(rows))

	for key, val := range rows {
		deliveryMethods[key] = entity.DeliveryMethod{
			ID:   val.ID,
			Name: val.Name,
			Slug: val.Slug,
		}
	}

	return deliveryMethods, nil
}

func (d PsqlDeliveryReadRepository) getDeliveryMethodBySlug(slug string) (*entity.DeliveryMethod, error) {

	var row DeliveryMethod

	err := d.db.Select("d.*").
		From(tableWithAlias(tableNameDeliveryMethods, "d")).
		Where(dbx.NewExp("d.slug={:slug}", dbx.Params{"slug": slug})).
		AndWhere(dbx.NewExp("d.status={:status}", dbx.Params{"status": true})).
		One(&row)

	if err != nil {
		return nil, err
	}

	return &entity.DeliveryMethod{
		ID:   row.ID,
		Name: row.Name,
		Slug: row.Slug,
	}, nil
}

func (d PsqlDeliveryReadRepository) getPaymentMethodsByDeliveryMethodId(ctx context.Context, id int) ([]entity.PaymentMethod, error) {

	var rows []PaymentMethod

	err := d.db.Select("p.*").
		From(tableWithAlias(tableNamePaymentMethods, "p")).
		InnerJoin(tableWithAlias(tableNameDeliveryAssignPayment, "a"), dbx.NewExp("p.id = a.payment_method_id")).
		Where(dbx.NewExp("a.delivery_method_id={:id}", dbx.Params{"id": id})).
		AndWhere(dbx.NewExp("p.status={:status}", dbx.Params{"status": true})).
		AndWhere(dbx.NewExp("a.status={:status}", dbx.Params{"status": true})).
		All(&rows)

	if err != nil {
		return nil, err
	}

	pMethods := make([]entity.PaymentMethod, len(rows))

	for k, v := range rows {
		pMethods[k] = entity.PaymentMethod{
			ID:   v.ID,
			Name: v.Name,
			Slug: v.Slug,
		}
	}

	return pMethods, nil
}

func (d PsqlDeliveryReadRepository) getPaymentMethodBySlug(slug string) (*entity.PaymentMethod, error) {

	var row PaymentMethod

	err := d.db.Select("p.*").
		From(tableWithAlias(tableNamePaymentMethods, "p")).
		Where(dbx.NewExp("p.slug={:slug}", dbx.Params{"slug": slug})).
		AndWhere(dbx.NewExp("p.status={:status}", dbx.Params{"status": true})).
		One(&row)

	if err != nil {
		return nil, err
	}

	return &entity.PaymentMethod{
		ID:   row.ID,
		Name: row.Name,
		Slug: row.Slug,
	}, nil
}

func (d PsqlDeliveryReadRepository) getWarehousesForYourselfByCity(ctx context.Context, city entity.City) ([]entity.Warehouse, error) {

	var row DeliveryAssignmentCity

	err := d.db.Select("dac.*").
		From(tableWithAlias(tableNameDeliveryAssignCity, "dac")).
		InnerJoin(tableWithAlias(tableNameDeliveryMethods, "d"), dbx.NewExp("dac.delivery_method_id = d.id")).
		Where(dbx.NewExp("d.slug={:slug}", dbx.Params{"slug": entity.DeliveryMethodYourself})).
		AndWhere(dbx.NewExp("dac.city_token={:token}", dbx.Params{"token": city.ID})).
		One(&row)

	if err != nil {
		return make([]entity.Warehouse, 0), err
	}

	if row.Warehouses.Valid == false {
		return make([]entity.Warehouse, 0), nil
	}

	var wj []Warehouse
	_ = json.Unmarshal([]byte(row.Warehouses.String), &wj)

	w := make([]entity.Warehouse, len(wj))

	for k, v := range wj {
		w[k] = entity.Warehouse{
			ID:      v.Name,
			Name:    v.Name,
			Address: v.Address,
			Phone:   v.Phone,
		}
	}

	return w, nil
}

func tableWithAlias(tableName, alias string) string {
	return tableName + " " + alias
}
