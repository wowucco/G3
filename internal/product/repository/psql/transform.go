package psql

import "github.com/wowucco/G3/internal/entity"

func rowToProductEntity(row *Product) *entity.Product {

	var (
		product entity.Product
		price entity.Price
		brand entity.Brand
		category entity.Category
		country entity.Country
		unit entity.Unit
		mainPhoto entity.Photo
	)

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
		Price: price,
	}

	return &product
}

func rowsToProductEntities(rows []Product) []*entity.Product {

	products := make([]*entity.Product, len(rows))

	for i, v := range rows {
		products[i] = rowToProductEntity(&v)
	}

	return products
}