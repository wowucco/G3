package psql

import (
	"context"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/entity"
	"strconv"
)

const tableNameProduct = "shop_products"

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
		product entity.Product
		price entity.Price
		brand entity.Brand
		category entity.Category
		country entity.Country
		unit entity.Unit
		mainPhoto entity.Photo
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

	if row.PhotoID.Valid == true {
		mainPhoto = entity.Photo{
			ID:   stringIdToInt(row.PhotoID.String),
			Link: row.PhotoFile.String,
			Main: true,
			Rating: stringIdToInt(row.PhotoSort.String),
		}
	}

	if row.BrandID.Valid == true {
		brand = entity.Brand{
			ID:   stringIdToInt(row.BrandID.String),
			Name: row.BrandName.String,
			Slug: row.BrandSlug.String,
		}
	}

	if row.CategoryID.Valid == true {
		category = entity.Category{
			ID:          stringIdToInt(row.CategoryID.String),
			Name:        row.CategoryName.String,
			Title:       row.CategoryTitle.String,
			Description: row.CategoryDescription.String,
			Slug:        row.CategorySlug.String,
		}
	}

	if row.CountryID.Valid == true {
		country = entity.Country{
			ID:   stringIdToInt(row.CountryID.String),
			Name: row.CountryName.String,
		}
	}

	if row.UnitID.Valid == true {
		unit = entity.Unit{
			ID:   stringIdToInt(row.UnitID.String),
			Name: row.UnitName.String,
		}
	}

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

	price = entity.Price{
		Price:     row.Price.Price,
		Currency: entity.Currency{
			ID:   row.CurrencyID,
			Name: row.CurrencyName,
			Rate: row.CurrencyRate,
			ISO:  row.CurrencyISO,
		},
	}

	if row.Price.SalePrice.Valid == true {
		price.SalePrice = stringIdToInt(row.Price.SalePrice.String)
		price.SaleCount = stringIdToInt(row.Price.SaleCount.String)
	}

	product = entity.Product{
		ID: row.ID,
		Name: row.Name,
		Description: row.Description,
		Code: row.Code,
		Exist: row.Exist,
		Status: row.Status,
		Group: entity.Group{
			ID: row.GroupID,
			Name: row.GroupName,
			Description: row.GroupDescription,
		},
		Country: country,
		Brand: brand,
		Category: category,
		Unit: unit,
		MainPhoto: mainPhoto,
		Photos: photos,
		Values: values,
		Price: price,
	}

	return &product, nil
}

func (r ProductRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.Select("COUNT(*)").From(tableNameProduct).Row(&count)
	return count, err
}

func (r ProductRepository) Query(ctx context.Context, offset, limit int) ([]*entity.Product, error) {
	//var products []*entity.Product
	//return products, nil

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

	products := make([]*entity.Product, len(rows))

	for i, v := range rows {

		var (
			price entity.Price
			brand entity.Brand
			category entity.Category
			country entity.Country
			unit entity.Unit
			mainPhoto entity.Photo
		)

		if v.PhotoID.Valid == true {
			mainPhoto = entity.Photo{
				ID:   stringIdToInt(v.PhotoID.String),
				Link: v.PhotoFile.String,
				Main: true,
				Rating: stringIdToInt(v.PhotoSort.String),
			}
		}

		if v.BrandID.Valid == true {
			brand = entity.Brand{
				ID:   stringIdToInt(v.BrandID.String),
				Name: v.BrandName.String,
				Slug: v.BrandSlug.String,
			}
		}

		if v.CategoryID.Valid == true {
			category = entity.Category{
				ID:          stringIdToInt(v.CategoryID.String),
				Name:        v.CategoryName.String,
				Title:       v.CategoryTitle.String,
				Description: v.CategoryDescription.String,
				Slug:        v.CategorySlug.String,
			}
		}

		if v.CountryID.Valid == true {
			country = entity.Country{
				ID:   stringIdToInt(v.CountryID.String),
				Name: v.CountryName.String,
			}
		}

		if v.UnitID.Valid == true {
			unit = entity.Unit{
				ID:   stringIdToInt(v.UnitID.String),
				Name: v.UnitName.String,
			}
		}

		price = entity.Price{
			Price:     v.Price.Price,
			Currency: entity.Currency{
				ID:   v.CurrencyID,
				Name: v.CurrencyName,
				Rate: v.CurrencyRate,
				ISO:  v.CurrencyISO,
			},
		}

		if v.Price.SalePrice.Valid == true {
			price.SalePrice = stringIdToInt(v.Price.SalePrice.String)
			price.SaleCount = stringIdToInt(v.Price.SaleCount.String)
		}

		products[i] = &entity.Product{
			ID: v.ID,
			Name: v.Name,
			Description: v.Description,
			Code: v.Code,
			Exist: v.Exist,
			Status: v.Status,
			Group: entity.Group{
				ID: v.GroupID,
				Name: v.GroupName,
				Description: v.GroupDescription,
			},
			Country: country,
			Brand: brand,
			Category: category,
			Unit: unit,
			MainPhoto: mainPhoto,
			Price: price,
		}
	}

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

func tableWithAlias(tableName, alias string) string {
	return tableName + " " + alias
}