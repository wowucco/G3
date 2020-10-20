package psql

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"strconv"
)

const tableNameProduct = "shop_products"
const tableNameBrands = "shop_brands"
const tableNameCategories = "shop_categories"
const tableNameCountry = "shop_country"
const tableNameProductUnit = "shop_products_unit"
const tableNamePhotos = "shop_photos"
const tableNameGroup = "shop_group"
const tableNameCurrency = "shop_currency"
const tableNameProductViewCount = "shop_product_view_count"
const tableNameOrderItems = "shop_order_items"

func stringIdToInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}

func (r ProductRepository) withValues(ids []interface{}) []CharacteristicValue {

	var values []CharacteristicValue

	_ = r.db.Select(
		"cv.id char_value_id", "cv.value char_value_value",
		"c.id char_id", "c.name char_name",
		"ct.id char_type_id", "ct.name char_type_name", "ct.custom char_type_custom",
		"u.id unit_id", "u.name unit_name",
		"v.product_id char_value_product_id").
		From("shop_values v").
		InnerJoin("shop_characteristics_values cv", dbx.NewExp("cv.id = v.value_id")).
		InnerJoin("shop_characteristics c", dbx.NewExp("v.characteristic_id = c.id")).
		InnerJoin("shop_characteristic_type ct", dbx.NewExp("c.type_id = ct.id")).
		LeftJoin("shop_products_unit u", dbx.NewExp("c.unit_id = u.id")).
		Where(dbx.In("v.product_id", ids...)).
		All(&values)

	return values
}

func (r ProductRepository) withPhotos(ids []interface{}) []Photo {

	var photos []Photo

	_ = r.db.Select("p.id photo_id", "p.file photo_file", "p.product_id photo_product_id", "p.sort photo_sort").
		From("shop_photos p").
		Where(dbx.In("p.product_id", ids...)).
		All(&photos)

	return photos
}

func (r ProductReadRepository) enabledCount(ctx context.Context) (int, error) {

	var count int

	err := r.db.Select("COUNT(*)").From(tableWithAlias(tableNameProduct, "p")).
		Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
		Row(&count)
	return count, err
}

func tableWithAlias(tableName, alias string) string {
	return tableName + " " + alias
}
