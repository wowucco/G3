package repository

import "database/sql"

type DeliveryMethod struct {
	ID     int            `db:"id"`
	Name   string         `db:"name"`
	Slug   string         `db:"slug"`
	Status bool           `db:"status"`
	Alpha  bool           `db:"alpha"`
	Width  sql.NullString `db:"width"`
}

type PaymentMethod struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Slug   string `db:"slug"`
	Status bool   `db:"status"`
}

type Warehouse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type NPWarehouse struct {
	ID                           string `json:"id"`
	CityID                       string `json:"cityId"`
	Name                         string `json:"name"`
	NameRu                       string `json:"nameRu"`
	ShortAddress                 string `json:"shortAddress"`
	ShortAddressRu               string `json:"shortAddressRu"`
	Phone                        string `json:"phone"`
	CityDescription              string `json:"cityDescription"`
	CityDescriptionRu            string `json:"cityDescriptionRu"`
	SettlementRef                string `json:"settlementRef"`
	SettlementDescription        string `json:"settlementDescription"`
	SettlementAreaDescription    string `json:"settlementAreaDescription"`
	SettlementRegionsDescription string `json:"settlementRegionsDescription"`
	SettlementTypeDescription    string `json:"settlementTypeDescription,omitempty"`
	Number                       int    `json:"number"`
	TotalMaxWeightAllowed        int    `json:"totalMaxWeightAllowed"`
}

type DeliveryAssignmentCity struct {
	Warehouses sql.NullString `db:"warehouses"`
}
