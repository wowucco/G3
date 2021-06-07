package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wowucco/G3/internal/entity"
	"log"
	"strconv"
	"strings"
)

const ESDeliveryIndex = "delivery"

const ESDeliveryCityDocType = "city"
const ESDeliveryWarehouseDocType = "warehouse"

func (d ESDeliveryReadRepository) getCityById(ctx context.Context, id string) (*entity.City, error) {

	q := map[string]interface{}{
		"_source": []string{
			"id",
			"nameRu",
		},
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"_id": []string{
					id,
				},
			},
		},
	}

	result, err := d.baseQueryToEs(ctx, q, ESDeliveryIndex, ESDeliveryCityDocType)

	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		return &entity.City{
			ID:   fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["id"]),
			Name: fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["nameRu"]),
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("city not found by %s", id))
}

func (d ESDeliveryReadRepository) searchCity(ctx context.Context, text string) ([]*entity.City, error) {

	input := strings.ToLower(text)
	should := make([]interface{}, 3)

	should[0] = map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": input,
			"fields": []string{
				"name",
				"nameRu",
			},
		},
	}
	should[1] = map[string]interface{}{
		"wildcard": map[string]string{
			"name": fmt.Sprintf("%s%s%s", "*", input, "*"),
		},
	}
	should[2] = map[string]interface{}{
		"wildcard": map[string]string{
			"nameRu": fmt.Sprintf("%s%s%s", "*", input, "*"),
		},
	}

	q := map[string]interface{}{
		"_source": []string{
			"id",
			"nameRu",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
			},
		},
	}

	result, err := d.baseQueryToEs(ctx, q, ESDeliveryIndex, ESDeliveryCityDocType)

	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	i := 0
	cities := make([]*entity.City, len(result["hits"].(map[string]interface{})["hits"].([]interface{})))

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		cities[i] = &entity.City{
			ID:   fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["id"]),
			Name: fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["nameRu"]),
		}
		i++
	}

	return cities, nil
}

func (d ESDeliveryReadRepository) baseQueryToEs(ctx context.Context, q map[string]interface{}, index, document string) (map[string]interface{}, error) {

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := d.es.Search(
		d.es.Search.WithContext(context.Background()),
		d.es.Search.WithIndex(index),
		d.es.Search.WithDocumentType(document),
		d.es.Search.WithBody(&buf),
		d.es.Search.WithPretty(),
	)

	defer res.Body.Close()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
			return nil, err
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)

			return nil, errors.New(fmt.Sprintf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"]))
		}
	}

	var result map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return nil, err
	}

	return result, nil
}

func (d ESDeliveryReadRepository) getWarehousesForNovaposhtaByCity(ctx context.Context, city entity.City) ([]NPWarehouse, error) {

	q := map[string]interface{}{
		"_source": []string{
			"id",
			"nameRu",
			"shortAddressRu",
			"phone",
			"number",
			"totalMaxWeightAllowed",
		},
		//"sort": []interface{}{
		//	map[string]interface{}{
		//		"id": []interface{}{
		//			map[string]string{
		//				"order": "asc",
		//			},
		//		},
		//	},
		//},
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"cityId": map[string]string{
					"query":            city.ID,
					"zero_terms_query": "all",
					"operator":         "and",
				},
			},
		},
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := d.es.Search(
		d.es.Search.WithContext(context.Background()),
		d.es.Search.WithIndex(ESDeliveryIndex),
		d.es.Search.WithDocumentType(ESDeliveryWarehouseDocType),
		d.es.Search.WithSize(1000),
		d.es.Search.WithBody(&buf),
		d.es.Search.WithPretty(),
	)

	defer res.Body.Close()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
			return nil, err
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)

			return nil, errors.New(fmt.Sprintf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"]))
		}
	}

	var result map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return nil, err
	}



	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	i := 0
	warehouses := make([]NPWarehouse, len(result["hits"].(map[string]interface{})["hits"].([]interface{})))

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {

		t, _ := strconv.Atoi(fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["totalMaxWeightAllowed"]))
		n, _ := strconv.Atoi(fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["number"]))

		warehouses[i] = NPWarehouse{
			ID:                           fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["id"]),
			CityID:                       fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["cityId"]),
			Name:                         fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["name"]),
			NameRu:                       fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["nameRu"]),
			ShortAddress:                 fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["shortAddress"]),
			ShortAddressRu:               fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["shortAddressRu"]),
			Phone:                        fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["phone"]),
			CityDescription:              fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["cityDescription"]),
			CityDescriptionRu:            fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["cityDescriptionRu"]),
			SettlementRef:                fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["settlementRef"]),
			SettlementDescription:        fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["settlementDescription"]),
			SettlementAreaDescription:    fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["settlementAreaDescription"]),
			SettlementRegionsDescription: fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["settlementRegionsDescription"]),
			SettlementTypeDescription:    fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["settlementTypeDescription"]),
			Number:                       n,
			TotalMaxWeightAllowed:        t,
		}

		log.Printf("Error getting response: %v\n", warehouses[i])

		i++
	}

	return warehouses, nil
}

func (d ESDeliveryReadRepository) indexCitiesWarehouses(city entity.City, warehouses []NPWarehouse) error {

	for _, warehouse := range warehouses {

		var buf bytes.Buffer

		if err := json.NewEncoder(&buf).Encode(warehouse); err != nil {
			log.Fatalf("Error encoding query: %s", err)
			return err
		}

		res, err := d.es.Index(ESDeliveryIndex, &buf, d.es.Index.WithDocumentID(warehouse.ID), d.es.Index.WithDocumentType(ESDeliveryWarehouseDocType))

		if err != nil {
			log.Fatalf("warhouse indexing get response error: %s", err)
		}

		res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document ID=%s", res.Status(), warehouse.ID)
		}
	}

	return nil
}

func (d ESDeliveryReadRepository) deleteWarehousesByCity(city entity.City) error {

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"cityId": map[string]string{
					"query":            city.ID,
					"zero_terms_query": "all",
					"operator":         "and",
				},
			},
		},
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return err
	}

	res, err := d.es.DeleteByQuery([]string{ESDeliveryIndex}, &buf, d.es.DeleteByQuery.WithDocumentType(ESDeliveryWarehouseDocType))
	defer res.Body.Close()

	return err
}

func (d ESDeliveryReadRepository) reindexCitiesWarehouses(city entity.City, warehouses []NPWarehouse) error {

	err := d.deleteWarehousesByCity(city)

	if err != nil {
		log.Fatalf("Error delete query: %s", err)
		return err
	}

	if len(warehouses) < 1 {
		return nil
	}

	err = d.indexCitiesWarehouses(city, warehouses)

	if err != nil {
		log.Fatalf("Error index query: %s", err)
		return err
	}

	return nil
}
