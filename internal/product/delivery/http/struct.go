package http

import "github.com/wowucco/G3/internal/entity"

type Photo struct {
	ID 		int		`json:"id"`
	Main 	bool	`json:"main"`
	Sort 	int		`json:"sort"`
	Small 	string  `json:"small"`
	Thumb 	string  `json:"thumb"`
}

type CharacteristicType struct {
	ID 		int 	`json:"id"`
	Name 	string 	`json:"name"`
	Custom 	bool 	`json:"custom"`
}

type Characteristic struct {
	ID 		int		`json:"id"`
	Name 	string	`json:"name"`
	Type 	CharacteristicType 	`json:"characteristic_type"`
	Unit 	Unit				`json:"unit"`
}

type CharacteristicValue struct {
	ID 				int		`json:"id"`
	Value 			string	`json:"value"`
	Characteristic 	Characteristic `json:"characteristic"`
}

type Country struct {
	ID 		int 	`json:"id,omitempty"`
	Name 	string 	`json:"name,omitempty"`
}

type Unit struct {
	ID 		int 	`json:"id,omitempty"`
	Name 	string 	`json:"name,omitempty"`
}

type Brand struct {
	ID 		int		`json:"id,omitempty"`
	Name 	string 	`json:"name,omitempty"`
	Slug 	string 	`json:"slug,omitempty"`
}

type Category struct {
	ID 			int		`json:"id,omitempty"`
	Name 		string	`json:"name,omitempty"`
	Title 		string	`json:"title,omitempty"`
	Description string	`json:"description,omitempty"`
	Slug 		string	`json:"slug,omitempty"`
}

type Group struct {
	ID 				int		`json:"id"`
	Name 			string	`json:"name"`
	Description 	string	`json:"description"`
}

type Price struct {
	Price 				string	`json:"price"`
	SalePrice 			string	`json:"sale_price,omitempty"`
	SaleCount 			int		`json:"sale_count,omitempty"`
	PriceInCents 		int		`json:"price_cents"`
	SalePriceInCents 	int		`json:"sale_price_cents,omitempty"`
	Currency  			string  `json:"currency"`
}

type Product struct {
	ID 			int		`json:"id"`
	Name 		string	`json:"name"`
	Description string	`json:"description"`
	Code 		int		`json:"code"`
	Exist 		int		`json:"exist"`
	Status 		int		`json:"status"`
	Price		Price	`json:"price"`
	Brand 		Brand		`json:"brand"`
	Category	Category 	`json:"category"`
	Group		Group		`json:"group"`
	Country		Country		`json:"country"`
	Unit		Unit		`json:"unit"`
	MainPhoto	Photo		`json:"main_photo"`
	Values		[]CharacteristicValue `json:"values"`
	Photos		[]Photo		`json:"photos"`
}

type getResponse struct {
	*Product `json:"product"`
}

func toProducts(p []*entity.Product) []*Product {
	products := make([]*Product, len(p))

	for i, v := range p {
		products[i] = toProduct(v)
	}

	return products
}

func toProduct(p *entity.Product) *Product {

	values := make([]CharacteristicValue, len(p.Values))

	for i, val := range p.Values {
		values[i] = CharacteristicValue{
			ID:             val.ID,
			Value:          val.Value,
			Characteristic: Characteristic{
				ID: val.Characteristic.ID,
				Name: val.Characteristic.Name,
				Type: CharacteristicType{
					ID:     val.Characteristic.Type.ID,
					Name:   val.Characteristic.Type.Name,
					Custom: val.Characteristic.Type.Custom,
				},
				Unit: Unit{
					ID:   val.Characteristic.Unit.ID,
					Name: val.Characteristic.Unit.Name,
				},
			},
		}
	}

	photos := make([]Photo, len(p.Photos))

	for i, val := range p.Photos {
		photos[i] = Photo{
			ID:   val.ID,
			Main: false,
			Sort: val.Rating,
			Small: val.GetSmallUrl(p),
			Thumb: val.GetThumbUrl(p),
		}
	}

	return &Product{
		ID:    			p.ID,
		Name: 			p.Name,
		Description: 	p.Description,
		Code: 			p.Code,
		Exist: 			p.Exist,
		Status: 		p.Status,
		Price: Price{
			Price:            p.Price.CentToCurrency(),
			SalePrice:        p.Price.SaleCentToCurrency(),
			SaleCount:        p.Price.SaleCount,
			PriceInCents:     p.Price.Price,
			SalePriceInCents: p.Price.SalePrice,
			Currency:         p.Price.Currency.Name,
		},
		Brand: Brand{
			ID: 	p.Brand.ID,
			Name: 	p.Brand.Name,
			Slug: 	p.Brand.Slug,
		},
		Group: Group{
			ID: 			p.Group.ID,
			Name: 			p.Group.Name,
			Description: 	p.Group.Description,
		},
		Category: Category{
			ID: p.Category.ID,
			Name: p.Category.Name,
			Title: p.Category.Title,
			Description: p.Category.Description,
			Slug: p.Category.Slug,
		},
		Country: Country{
			p.Country.ID,
			p.Country.Name,
		},
		Unit: Unit{
			p.Unit.ID,
			p.Unit.Name,
		},
		MainPhoto: Photo{
			ID:    p.MainPhoto.ID,
			Main:  p.MainPhoto.IsMain(),
			Sort:  p.MainPhoto.Rating,
			Small: p.MainPhoto.GetSmallUrl(p),
			Thumb: p.MainPhoto.GetThumbUrl(p),
		},
		Values: values,
		Photos: photos,
	}
}