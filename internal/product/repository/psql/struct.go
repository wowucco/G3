package psql

import (
	"database/sql"
)

type Meta struct {

}
type Photo struct {
	ID 			int 	`db:"photo_id"`
	File 		string 	`db:"photo_file"`
	ProductID 	int 	`db:"photo_product_id"`
	Sort 		int 	`db:"photo_sort"`
}

type MainPhoto struct {
	PhotoID 		sql.NullString 	`db:"photo_id"`
	PhotoFile 		sql.NullString 	`db:"photo_file"`
	PhotoProductID 	sql.NullString 	`db:"photo_product_id"`
	PhotoSort 		sql.NullString 	`db:"photo_sort"`
}

type Unit struct {
	UnitID 		sql.NullString `db:"unit_id"`
	UnitName 	sql.NullString `db:"unit_name"`
}

type CharacteristicType struct {
	CharacteristicTypeID 		int				`db:"char_type_id"`
	CharacteristicTypeName 		string			`db:"char_type_name"`
	CharacteristicTypeCustom 	sql.NullString	`db:"char_type_custom"`
}

type Characteristic struct {
	CharacteristicID 	int			`db:"char_id"`
	CharacteristicName 	string		`db:"char_name"`
	CharacteristicType
	Unit
}

type CharacteristicValue struct {
	ID 				int		`db:"char_value_id"`
	Value 			string	`db:"char_value_value"`
	ProductID 		int		`db:"char_value_product_id"`
	Characteristic
}

type Country struct {
	CountryID 	sql.NullString 	`db:"country_id"`
	CountryName sql.NullString	`db:"country_name"`
}

type Category struct {
	CategoryID 			sql.NullString	`db:"category_id"`
	CategoryName 		sql.NullString	`db:"category_name"`
	CategoryTitle 		sql.NullString	`db:"category_title"`
	CategoryDescription sql.NullString	`db:"category_description"`
	CategorySlug 		sql.NullString	`db:"category_slug"`
}

type Brand struct {
	BrandID 	sql.NullString	`db:"brand_id"`
	BrandName 	sql.NullString	`db:"brand_name"`
	BrandSlug 	sql.NullString	`db:"brand_slug"`
}

type Group struct {
	GroupID 			int		`db:"group_id"`
	GroupName 			string	`db:"group_name"`
	GroupDescription 	string	`db:"group_description"`
}

type Currency struct {
	CurrencyID 		int		`db:"currency_id"`
	CurrencyName 	string	`db:"currency_name"`
	CurrencyRate 	float32	`db:"currency_rate"`
	CurrencyISO 	string	`db:"currency_iso"`
}

type Price struct {
	Price 		int				`db:"price"`
	SalePrice 	sql.NullString	`db:"sale_price"`
	SaleCount 	sql.NullString	`db:"sale_count"`
	Currency
}

type Product struct {
	ID 			int
	Name 		string
	Description string
	Code 		int
	Exist 		int
	Status 		int
	Brand
	Group
	Category
	Country
	Unit
	Price
	MainPhoto
	Photos []Photo
}