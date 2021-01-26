package repository

import "database/sql"

type DeliveryMethod struct {
	ID   int 	`db:"id"`
	Name string `db:"name"`
	Slug string `db:"slug"`
	Status bool `db:"status"`
	Alpha bool `db:"alpha"`
	Width sql.NullString `db:"width"`
}

type PaymentMethod struct {
	ID   int 	`db:"id"`
	Name string `db:"name"`
	Slug string `db:"slug"`
	Status bool `db:"status"`
}

type Warehouse struct {
	Name string	`json:"name"`
	Address string `json:"address"`
	Phone string `json:"phone"`
}

type NPWarehouse struct {
	ID                           string
	CityID                       string
	Name                         string
	NameRu                       string
	ShortAddress                 string
	ShortAddressRu               string
	Phone                        string
	CityDescription              string
	CityDescriptionRu            string
	SettlementRef                string
	SettlementDescription        string
	SettlementAreaDescription    string
	SettlementRegionsDescription string
	SettlementTypeDescription    string
}

type DeliveryAssignmentCity struct {
	Warehouses sql.NullString `db:"warehouses"`
}
