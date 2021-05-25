package repository

import "database/sql"

type NextId struct {
	Id int
}

type City struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
type Address struct {
	IsCustom bool   `json:"is_custom"`
	Address  string `json:"address"`
	Code     string `json:"code"`
}
type PaymentExtra struct {
	Edrpou   string `json:"edrpou"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	PartsPay int    `json:"parts_pay"`
}
type Payment struct {
	ID            int            `db:"int"`
	TransactionID string         `db:"transaction_id"`
	OrderId       int            `db:"order_id"`
	Provider      string         `db:"provider"`
	Amount        int            `db:"amount"`
	Status        int            `db:"status"`
	Meta          sql.NullString `db:"meta"`
	Created       sql.NullString `db:"created_at"`
	Updated       sql.NullString `db:"updated_at"`
}
type DeliveryInfo struct {
	City    City         `json:"city"`
	Address Address      `json:"address"`
	Payment PaymentExtra `json:"payment"`
}

type DeliveryStatus struct {
	Status  int    `json:"status"`
	Time    int64  `json:"time"`
	Comment string `json:"comment"`
}

type PaymentStatus struct {
	Status  int    `json:"status"`
	Time    int64  `json:"time"`
	Comment string `json:"comment"`
}

type DeliveryMethod struct {
	DeliveryId   int    `db:"delivery_id"`
	DeliveryName string `db:"delivery_name"`
	DeliverySlug string `db:"delivery_slug"`
}
type PaymentMethod struct {
	PaymentId   int    `db:"payment_id"`
	PaymentName string `db:"payment_name"`
	PaymentSlug string `db:"payment_slug"`
}
type Customer struct {
	Name  string `db:"customer_name"`
	Phone string `db:"customer_phone"`
}
type Order struct {
	Id        int    `db:"id"`
	Created   int64  `db:"created_at"`
	Comment   string `db:"comment"`
	DoNotCall bool   `db:"do_not_call"`
	Cost      int    `db:"cost"`

	DeliveryStatus       int    `db:"delivery_status"`
	DeliveryStatusesJson string `db:"delivery_statuses_json"`
	DeliveryInfo         string `db:"delivery_info"`

	PaymentStatus       int    `db:"payment_status"`
	PaymentStatusesJson string `db:"payment_statuses_json"`

	Customer
	DeliveryMethod
	PaymentMethod
}
type Product struct {
	ProductId     int    `db:"product_id"`
	ProductName   string `db:"product_name"`
	ProductCode   int    `db:"product_code"`
	ProductExist  int    `db:"product_exist"`
	ProductStatus int    `db:"product_status"`

	ProductPrice     int            `db:"product_price"`
	ProductSalePrice sql.NullString `db:"product_sale_price"`
	ProductSaleCount sql.NullString `db:"product_sale_count"`
}
type Currency struct {
	CurrencyID   int     `db:"currency_id"`
	CurrencyName string  `db:"currency_name"`
	CurrencyRate float32 `db:"currency_rate"`
	CurrencyISO  string  `db:"currency_iso"`
}
type OrderItem struct {
	ID       int `db:"id"`
	Quantity int `db:"quantity"`
	Cost     int `db:"price"`
	Product
	Currency
}
