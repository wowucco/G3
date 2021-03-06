// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Category struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Title        string  `json:"title"`
	Descriptinon *string `json:"descriptinon"`
	Slug         string  `json:"slug"`
}

type Characteristic struct {
	ID   int                 `json:"id"`
	Name string              `json:"name"`
	Type *CharacteristicType `json:"type"`
	Unit *Unit               `json:"unit"`
}

type CharacteristicType struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsCustom bool   `json:"isCustom"`
}

type CharacteristicValue struct {
	ID             int             `json:"id"`
	Value          string          `json:"value"`
	Characteristic *Characteristic `json:"characteristic"`
}

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DeliveryInfo struct {
	DeliveryMethod *DeliveryMethod  `json:"deliveryMethod"`
	PaymentMethods []*PaymentMethod `json:"paymentMethods"`
	Warehouses     []*Warehouse     `json:"warehouses"`
}

type DeliveryMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ExistProduct struct {
	Exist bool `json:"exist"`
	ID    int  `json:"id"`
}

type Group struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type Pages struct {
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
	PageCount  int        `json:"pageCount"`
	TotalCount int        `json:"totalCount"`
	Items      []*Product `json:"items"`
}

type PagesWithGroup struct {
	Pages *Pages `json:"pages"`
	Group *Group `json:"group"`
}

type PagesWithGroups struct {
	Pages  *Pages   `json:"pages"`
	Groups []*Group `json:"groups"`
}

type PaymentMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Photo struct {
	ID     int    `json:"id"`
	IsMain bool   `json:"isMain"`
	Sort   int    `json:"sort"`
	Small  string `json:"small"`
	Thumb  string `json:"thumb"`
}

type Price struct {
	Price            string  `json:"price"`
	SalePrice        *string `json:"salePrice"`
	SaleCount        *int    `json:"saleCount"`
	PriceInCents     int     `json:"priceInCents"`
	SalePriceInCents *int    `json:"salePriceInCents"`
	Currency         string  `json:"currency"`
}

type Product struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description"`
	Code        int                    `json:"code"`
	Exist       int                    `json:"exist"`
	Status      int                    `json:"status"`
	Price       *Price                 `json:"price"`
	Brand       *Brand                 `json:"brand"`
	Category    *Category              `json:"category"`
	Group       *Group                 `json:"group"`
	Country     *Country               `json:"country"`
	Unit        *Unit                  `json:"unit"`
	MainPhoto   *Photo                 `json:"mainPhoto"`
	Photos      []*Photo               `json:"photos"`
	Values      []*CharacteristicValue `json:"values"`
}

type SimpleProduct struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Code        int       `json:"code"`
	Exist       int       `json:"exist"`
	Status      int       `json:"status"`
	Price       *Price    `json:"price"`
	Brand       *Brand    `json:"brand"`
	Category    *Category `json:"category"`
	Group       *Group    `json:"group"`
	Country     *Country  `json:"country"`
	Unit        *Unit     `json:"unit"`
	MainPhoto   *Photo    `json:"mainPhoto"`
}

type TreeChildrenMenuItem struct {
	ID          int                     `json:"id"`
	Name        string                  `json:"name"`
	Image       *string                 `json:"image"`
	Children    []*TreeChildrenMenuItem `json:"children"`
	HasChildren bool                    `json:"hasChildren"`
}

type TreeMenu struct {
	ID     *int  `json:"id"`
	Depht  *int  `json:"depht"`
	Parent *bool `json:"parent"`
}

type TreeMenuItem struct {
	ID          int                     `json:"id"`
	Name        string                  `json:"name"`
	Description *string                 `json:"description"`
	Image       *string                 `json:"image"`
	Parent      *TreeParentMenuItem     `json:"parent"`
	Children    []*TreeChildrenMenuItem `json:"children"`
	HasParent   bool                    `json:"hasParent"`
	HasChildren bool                    `json:"hasChildren"`
}

type TreeParentMenuItem struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Image     *string             `json:"image"`
	Parent    *TreeParentMenuItem `json:"parent"`
	HasParent bool                `json:"hasParent"`
}

type Unit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Warehouse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Number    int    `json:"number"`
	MaxWeight int    `json:"maxWeight"`
}

type CityID struct {
	ID string `json:"id"`
}

type ID struct {
	ID int `json:"id"`
}

type IDWithLimit struct {
	ID    int `json:"id"`
	Limit int `json:"limit"`
}

type Ids struct {
	Ids []*int `json:"ids"`
}

type Page struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type PageByID struct {
	ID      int `json:"id"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type PageByIds struct {
	Ids     []*int `json:"ids"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

type Text struct {
	Text string `json:"text"`
}
