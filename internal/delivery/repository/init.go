package repository

import (
	"github.com/elastic/go-elasticsearch/v5"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/go-novaposhta"
)

type DeliveryReadRepository struct {

	db *PsqlDeliveryReadRepository
	es *ESDeliveryReadRepository
	np *NPDeliveryReadRepository
}

type ESDeliveryReadRepository struct {
	es *elasticsearch.Client
}

type NPDeliveryReadRepository struct {
	np *novaposhta.Client
}

type PsqlDeliveryReadRepository struct {
	db *dbx.DB
}

func NewDeliveryReadRepository(db *dbx.DB, es *elasticsearch.Client, np *novaposhta.Client) *DeliveryReadRepository {
	return &DeliveryReadRepository{
		db: &PsqlDeliveryReadRepository{db:db},
		es: &ESDeliveryReadRepository{es:es},
		np: &NPDeliveryReadRepository{np:np},
	}
}
