package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strconv"
)

const baseCurrency = "UAH"

const photoLinkTypeOrigin = "origin"
const photoLinkTypeThumb = "thumb"
const photoLinkTypeSmall = "small"

type Photo struct {
	ID 		int
	Link 	string
	Main 	bool
	Rating 	int
}

func (photo Photo) IsMain() bool  {
	return photo.Main
}

func (photo *Photo) GetOriginUrl(p *Product) string {
	return photo.getUrl(p, photoLinkTypeOrigin)
}

func (photo *Photo) GetSmallUrl(p *Product) string {
	return photo.getUrl(p, photoLinkTypeSmall)
}

func (photo *Photo) GetThumbUrl(p *Product) string {
	return photo.getUrl(p, photoLinkTypeThumb)
}

func (photo *Photo) getUrl(p *Product, linkType string) string {

	h := md5.New()
	h.Write([]byte(strconv.Itoa(photo.ID)))
	hash := hex.EncodeToString(h.Sum(nil))

	if linkType == photoLinkTypeOrigin {
		return viper.GetString("domain.static") +
			"/media/products/" +
			strconv.Itoa(p.ID) +
			linkType +
			hash +
			filepath.Ext(photo.Link)
	}

	return viper.GetString("domain.static") +
		"/media/products/" +
		strconv.Itoa(p.ID) +
		"/cache/" + linkType + "_" +
		hash +
		filepath.Ext(photo.Link)
}

type Unit struct {
	ID 		int
	Name 	string
}

type CharacteristicType struct {
	ID 		int
	Name 	string
	Custom 	bool
}

type Characteristic struct {
	ID 		int
	Name 	string
	Type 	CharacteristicType
	Unit 	Unit
}

type CharacteristicValue struct {
	ID 				int
	Characteristic 	Characteristic
	Value 			string
}

type Country struct {
	ID 		int
	Name 	string
}

type Currency struct {
	ID 		int
	Name 	string
	Rate 	float32
	ISO 	string
}

func (c *Currency) IsBase() bool {
	if c.ISO == baseCurrency {
		return true
	}

	return false
}

type Price struct {
	Price 		int
	SalePrice 	int
	SaleCount 	int
	Currency 	Currency
}

func (p *Price) CentToCurrency() string {
	return p.toCurrency(p.Price)
}

func (p *Price) SaleCentToCurrency() string {

	if p.SalePrice == 0 {
		return ""
	}

	return p.toCurrency(p.SalePrice)
}

func (p *Price) toCurrency(cents int) string {
	return fmt.Sprintf("%.2f", float64(cents / 100))
}

type Meta struct {
	Keywords 	string
	Title 		string
	Description string
}

type Brand struct {
	ID 		int
	Name 	string
	Slug 	string
}

type Category struct {
	ID 			int
	Name 		string
	Title		string
	Description string
	Slug 		string
	Photo 		Photo
	Meta 		Meta
}

type Group struct {
	ID 			int
	Name 		string
	Description string
	Photo 		Photo
	Meta 		Meta
}

type Product struct {
	ID 			int
	Name 		string
	Description string
	Code 		int
	Exist 		int
	Status 		int

	Brand 		Brand
	Category 	Category
	Group		Group
	Unit		Unit
	Country 	Country

	Price 		Price

	Values []CharacteristicValue

	MainPhoto 	Photo
	Photos 		[]Photo

	Meta 		Meta
}