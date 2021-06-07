package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wowucco/G3/internal/entity"
	"log"
	"strconv"
)

func (d NPDeliveryReadRepository) getWarehousesForNovaposhtaByCity(ctx context.Context, city entity.City) ([]NPWarehouse, error) {

	res, _ := d.np.Address.GetWarehouses(d.np.Address.GetWarehouses.WithParams("", city.ID, "ua", 0, 0))

	defer res.Body.Close()

	var (
		result map[string]interface{}
	)
	//
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("error %v\n", err)
	}

	if result["success"] != true {
		err := fmt.Errorf("KO response from novaposhta\n")
		log.Printf("%v\n", err)
		return nil, err
	}

	w := make([]NPWarehouse, len(result["data"].([]interface{})))

	for k, warehouse := range result["data"].([]interface{}) {

		t, _ := strconv.Atoi(fmt.Sprintf("%v", warehouse.(map[string]interface{})["totalMaxWeightAllowed"]))
		n, _ := strconv.Atoi(fmt.Sprintf("%v", warehouse.(map[string]interface{})["number"]))

		w[k] = NPWarehouse{
			ID:                           fmt.Sprintf("%v", warehouse.(map[string]interface{})["Ref"]),
			CityID:                       fmt.Sprintf("%v", warehouse.(map[string]interface{})["CityRef"]),
			Name:                         fmt.Sprintf("%v", warehouse.(map[string]interface{})["Description"]),
			NameRu:                       fmt.Sprintf("%v", warehouse.(map[string]interface{})["DescriptionRu"]),
			ShortAddress:                 fmt.Sprintf("%v", warehouse.(map[string]interface{})["ShortAddress"]),
			ShortAddressRu:               fmt.Sprintf("%v", warehouse.(map[string]interface{})["ShortAddressRu"]),
			Phone:                        fmt.Sprintf("%v", warehouse.(map[string]interface{})["Phone"]),
			CityDescription:              fmt.Sprintf("%v", warehouse.(map[string]interface{})["CityDescription"]),
			CityDescriptionRu:            fmt.Sprintf("%v", warehouse.(map[string]interface{})["CityDescriptionRu"]),
			SettlementRef:                fmt.Sprintf("%v", warehouse.(map[string]interface{})["SettlementRef"]),
			SettlementDescription:        fmt.Sprintf("%v", warehouse.(map[string]interface{})["SettlementDescription"]),
			SettlementAreaDescription:    fmt.Sprintf("%v", warehouse.(map[string]interface{})["SettlementAreaDescription"]),
			SettlementRegionsDescription: fmt.Sprintf("%v", warehouse.(map[string]interface{})["SettlementRegionsDescription"]),
			SettlementTypeDescription:    fmt.Sprintf("%v", warehouse.(map[string]interface{})["SettlementTypeDescription"]),
			Number:                       n,
			TotalMaxWeightAllowed:        t,
		}
	}

	return w, nil
}
