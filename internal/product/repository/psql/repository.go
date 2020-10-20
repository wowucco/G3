package psql

import (
	"context"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/entity"
)

type ProductRepository struct {
	db *dbx.DB
}

func NewProductRepository(db *dbx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r ProductRepository) Get(ctx context.Context, id int) (*entity.Product, error) {

	var (
		row Product
		product *entity.Product
	)

	err := r.db.Select(
		"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
			"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
			"g.id group_id", "g.name group_name", "g.description group_description",
			"cnt.id country_id", "cnt.name country_name",
			"u.id unit_id", "u.name unit_name",
			"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
			"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
			"p.*").
		From(tableWithAlias(tableNameProduct, "p")).
			LeftJoin("shop_brands b", dbx.NewExp("p.brand_id = b.id")).
			LeftJoin("shop_categories c", dbx.NewExp("p.category_id = c.id")).
			LeftJoin("shop_country cnt", dbx.NewExp("p.country_id = cnt.id")).
			LeftJoin("shop_products_unit u", dbx.NewExp("p.unit_id = u.id")).
			LeftJoin("shop_photos ph", dbx.NewExp("p.main_photo_id = ph.id")).
			InnerJoin("shop_group g", dbx.NewExp("p.group_id = g.id")).
			InnerJoin("shop_currency cr", dbx.NewExp("p.currency_id = cr.id")).
		Where(dbx.NewExp("p.id={:id}", dbx.Params{"id": id})).
		One(&row)

	if err != nil {
		fmt.Print(err.Error())
	}

	product = rowToProductEntity(&row)

	ids := make([]interface{}, 1)
	ids[0] = id

	wPhotos := r.withPhotos(ids)
	photos := make([]entity.Photo, len(wPhotos))

	for i, val := range wPhotos {
		photos[i] = entity.Photo{
			ID:     val.ID,
			Link:   val.File,
			Main:   false,
			Rating: val.Sort,
		}
	}

	product.Photos = photos

	wValues := r.withValues(ids)
	values := make([]entity.CharacteristicValue, len(wValues))

	for i, val := range wValues {

		unit := entity.Unit{}
		if val.UnitID.Valid == true {
			unit.ID = stringIdToInt(val.UnitID.String)
			unit.Name = val.UnitName.String
		}

		values[i] = entity.CharacteristicValue{
			ID: val.ID,
			Value: val.Value,
			Characteristic: entity.Characteristic{
				ID:   val.CharacteristicID,
				Name: val.CharacteristicName,
				Type: entity.CharacteristicType{
					ID: val.CharacteristicTypeID,
					Name: val.CharacteristicTypeName,
				},
				Unit: unit,
			},
		}
	}

	product.Values = values

	return product, nil
}

func (r ProductRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.Select("COUNT(*)").From(tableNameProduct).Row(&count)
	return count, err
}

func (r ProductRepository) Query(ctx context.Context, offset, limit int) ([]*entity.Product, error) {

	var (
		rows []Product
	)

	err := r.db.Select(
			"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
			"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
			"g.id group_id", "g.name group_name", "g.description group_description",
			"cnt.id country_id", "cnt.name country_name",
			"u.id unit_id", "u.name unit_name",
			"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
			"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
			"p.*").
		From(tableWithAlias(tableNameProduct, "p")).
			LeftJoin("shop_brands b", dbx.NewExp("p.brand_id = b.id")).
			LeftJoin("shop_categories c", dbx.NewExp("p.category_id = c.id")).
			LeftJoin("shop_country cnt", dbx.NewExp("p.country_id = cnt.id")).
			LeftJoin("shop_products_unit u", dbx.NewExp("p.unit_id = u.id")).
			LeftJoin("shop_photos ph", dbx.NewExp("p.main_photo_id = ph.id")).
			InnerJoin("shop_group g", dbx.NewExp("p.group_id = g.id")).
			InnerJoin("shop_currency cr", dbx.NewExp("p.currency_id = cr.id")).
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&rows)

	products := rowsToProductEntities(rows)

	return products, err
}

func (r ProductRepository) Create(ctx context.Context, product *entity.Product) (int, error) {
	return 0, nil
}

func (r ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	return nil
}

func (r ProductRepository) Delete(ctx context.Context, id int) error {
	return nil
}