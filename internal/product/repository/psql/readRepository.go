package psql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v5"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/entity"
	"log"
	"strings"
)

type ProductReadRepository struct {
	db *dbx.DB
	es *elasticsearch.Client
}

func NewProductReadRepository(db *dbx.DB, es *elasticsearch.Client) *ProductReadRepository {
	return &ProductReadRepository{
		db: db,
		es: es,
	}
}

func (r ProductReadRepository) GetById(ctx context.Context, id int, with []string) (*entity.Product, error) {

	var (
		row Product
		product *entity.Product
	)

	err := r.db.Select(
			"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
			"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
			"g.id group_id", "g.name group_name", "g.description group_description",
			"cnt.id country_id", "cnt.name country_name",
			"u.id unit_id", "u.name unit_name",
			"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
			"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
			"p.*").
		From(tableWithAlias(tableNameProduct, "p")).
			LeftJoin("shop_brands b", dbx.NewExp("p.brand_id = b.id")).
			LeftJoin("shop_categories c", dbx.NewExp("p.category_id = c.id")).
			LeftJoin("shop_country cnt", dbx.NewExp("p.country_id = cnt.id")).
			LeftJoin("shop_products_unit u", dbx.NewExp("p.unit_id = u.id")).
			LeftJoin("shop_photos ph", dbx.NewExp("p.main_photo_id = ph.id")).
			InnerJoin("shop_group g", dbx.NewExp("p.group_id = g.id")).
			InnerJoin("shop_currency cr", dbx.NewExp("p.currency_id = cr.id")).
		Where(dbx.NewExp("p.id={:id}", dbx.Params{"id": id})).
		One(&row)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	product = rowToProductEntity(&row)

	if product.ID == 0 {
		return nil, fmt.Errorf("product %d not found", id)
	}

	return product, nil
}

func (r ProductReadRepository) GetByIdsWithSequence(ctx context.Context, ids []int) ([]*entity.Product, error) {

	var (
		rows []Product
		sequence string
	)

	pIds := make([]interface{}, len(ids))

	for k, v := range ids {
		pIds[k] = v

		if sequence != "" {
			sequence += ","
		}

		sequence += fmt.Sprintf("(%d, %d)", k, v)
	}

	q := r.db.Select(
		"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
		"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
		"g.id group_id", "g.name group_name", "g.description group_description",
		"cnt.id country_id", "cnt.name country_name",
		"u.id unit_id", "u.name unit_name",
		"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
		"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
		"p.*").
	From(tableWithAlias(tableNameProduct, "p")).
		LeftJoin(tableWithAlias(tableNameBrands, "b"), dbx.NewExp("p.brand_id = b.id")).
		LeftJoin(tableWithAlias(tableNameCategories, "c"), dbx.NewExp("p.category_id = c.id")).
		LeftJoin(tableWithAlias(tableNameCountry, "cnt"), dbx.NewExp("p.country_id = cnt.id")).
		LeftJoin(tableWithAlias(tableNameProductUnit, "u"), dbx.NewExp("p.unit_id = u.id")).
		LeftJoin(tableWithAlias(tableNamePhotos, "ph"), dbx.NewExp("p.main_photo_id = ph.id")).
		LeftJoin(tableWithAlias(tableNameProductViewCount, "vc"), dbx.NewExp("p.id = vc.product_id")).
		InnerJoin(tableWithAlias(tableNameGroup, "g"), dbx.NewExp("p.group_id = g.id")).
		InnerJoin(tableWithAlias(tableNameCurrency, "cr"), dbx.NewExp("p.currency_id = cr.id")).
	Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
	Where(dbx.In("p.id", pIds...))

	if sequence != "" {
		q.InnerJoin(fmt.Sprintf("(values %s ) as last (ordering, id)", sequence), dbx.NewExp("p.id = last.id")).
			OrderBy("last.ordering")
	}

	err := q.All(&rows)

	return rowsToProductEntities(rows), err
}

func (r ProductReadRepository) GetSimilar(ctx context.Context, p entity.Product, size int) ([]*entity.Product, error) {

	var (
		buf bytes.Buffer
		split []string
	)

	name := strings.ToLower(p.Name)

	split = strings.Split(name, " ")

	should := make([]interface{}, len(split) + 1)

	should[0] = map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": name,
			"fields": []string{
				"name",
			},
		},
	}

	for k, v := range split {
		should[k+1] = map[string]interface{}{
			"wildcard": map[string]string{
				"name": "*" + v + "*",
			},
		}
	}

	q := map[string]interface{}{
		"_source": []string{
			"id",
		},
		"size": size,
		"sort": []interface{}{
			map[string]interface{}{
				"_score": []interface{}{
					map[string]string{
						"order": "desc",
					},
				},
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
				"must_not": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"id": p.ID,
						},
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("shop"),
		r.es.Search.WithDocumentType("products"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithPretty(),
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

	var (
		result map[string]interface{}
	)

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return nil, err
	}

	i := 0
	ids := make([]int, len(result["hits"].(map[string]interface{})["hits"].([]interface{})))

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		ids[i] = stringIdToInt(fmt.Sprintf("%v", hit.(map[string]interface{})["_id"]))
		i++
	}

	return r.GetByIdsWithSequence(ctx, ids)
}

func (r ProductReadRepository) GetPopularCount(ctx context.Context) (int, error) {

	return r.enabledCount(ctx)
}

func (r ProductReadRepository) GetPopular(ctx context.Context, offset, limit int) ([]*entity.Product, error) {

	var rows []Product

	err := r.db.Select(
		"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
		"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
		"g.id group_id", "g.name group_name", "g.description group_description",
		"cnt.id country_id", "cnt.name country_name",
		"u.id unit_id", "u.name unit_name",
		"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
		"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
		"p.*").
	From(tableWithAlias(tableNameProduct, "p")).
		LeftJoin(tableWithAlias(tableNameBrands, "b"), dbx.NewExp("p.brand_id = b.id")).
		LeftJoin(tableWithAlias(tableNameCategories, "c"), dbx.NewExp("p.category_id = c.id")).
		LeftJoin(tableWithAlias(tableNameCountry, "cnt"), dbx.NewExp("p.country_id = cnt.id")).
		LeftJoin(tableWithAlias(tableNameProductUnit, "u"), dbx.NewExp("p.unit_id = u.id")).
		LeftJoin(tableWithAlias(tableNamePhotos, "ph"), dbx.NewExp("p.main_photo_id = ph.id")).
		LeftJoin(tableWithAlias(tableNameProductViewCount, "vc"), dbx.NewExp("p.id = vc.product_id")).
		InnerJoin(tableWithAlias(tableNameGroup, "g"), dbx.NewExp("p.group_id = g.id")).
		InnerJoin(tableWithAlias(tableNameCurrency, "cr"), dbx.NewExp("p.currency_id = cr.id")).
	Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
	OrderBy("case when \"vc\".\"count\" is null then 1 else 0 end, \"vc\".\"count\" desc").
	Offset(int64(offset)).
	Limit(int64(limit)).
	All(&rows)

	return rowsToProductEntities(rows), err
}

func (r ProductReadRepository) GetTopSalesCount(ctx context.Context) (int, error) {

	var count int

	err := r.db.Select("COUNT(*)").From(tableWithAlias(tableNameProduct, "p")).
		InnerJoin(tableWithAlias(tableNameOrderItems, "oi"), dbx.NewExp("p.id = oi.product_id")).
		Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
		Row(&count)

	return count, err
}

func (r ProductReadRepository) GetTopSales(ctx context.Context, offset, limit int) ([]*entity.Product, error) {

	var rows []Product

	err := r.db.Select(
		"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
		"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
		"g.id group_id", "g.name group_name", "g.description group_description",
		"cnt.id country_id", "cnt.name country_name",
		"u.id unit_id", "u.name unit_name",
		"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
		"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
		"sum(oi.quantity) summa",
		"p.*").
	From(tableWithAlias(tableNameProduct, "p")).
		LeftJoin(tableWithAlias(tableNameBrands, "b"), dbx.NewExp("p.brand_id = b.id")).
		LeftJoin(tableWithAlias(tableNameCategories, "c"), dbx.NewExp("p.category_id = c.id")).
		LeftJoin(tableWithAlias(tableNameCountry, "cnt"), dbx.NewExp("p.country_id = cnt.id")).
		LeftJoin(tableWithAlias(tableNameProductUnit, "u"), dbx.NewExp("p.unit_id = u.id")).
		LeftJoin(tableWithAlias(tableNamePhotos, "ph"), dbx.NewExp("p.main_photo_id = ph.id")).
		InnerJoin(tableWithAlias(tableNameOrderItems, "oi"), dbx.NewExp("p.id = oi.product_id")).
		InnerJoin(tableWithAlias(tableNameGroup, "g"), dbx.NewExp("p.group_id = g.id")).
		InnerJoin(tableWithAlias(tableNameCurrency, "cr"), dbx.NewExp("p.currency_id = cr.id")).
	Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
	GroupBy("p.id", "b.id", "c.id", "g.id", "cnt.id", "u.id", "ph.id", "cr.id").
	OrderBy("summa desc").
	Offset(int64(offset)).
	Limit(int64(limit)).
	All(&rows)

	return rowsToProductEntities(rows), err
}

func (r ProductReadRepository) GetPopularByGroupIdCount(ctx context.Context, groupId int) (int, error) {

	var count int

	err := r.db.Select("COUNT(*)").From(tableWithAlias(tableNameProduct, "p")).
		Where(dbx.NewExp("p.status={:status} and p.group_id={:group_id}", dbx.Params{"status": 1, "group_id": groupId})).
		Row(&count)

	return count, err
}

func (r ProductReadRepository) GetPopularByGroupId(ctx context.Context, groupId int, offset, limit int) ([]*entity.Product, error) {
	var rows []Product

	err := r.db.Select(
			"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
			"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
			"g.id group_id", "g.name group_name", "g.description group_description",
			"cnt.id country_id", "cnt.name country_name",
			"u.id unit_id", "u.name unit_name",
			"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
			"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
			"p.*").
		From(tableWithAlias(tableNameProduct, "p")).
			LeftJoin(tableWithAlias(tableNameBrands, "b"), dbx.NewExp("p.brand_id = b.id")).
			LeftJoin(tableWithAlias(tableNameCategories, "c"), dbx.NewExp("p.category_id = c.id")).
			LeftJoin(tableWithAlias(tableNameCountry, "cnt"), dbx.NewExp("p.country_id = cnt.id")).
			LeftJoin(tableWithAlias(tableNameProductUnit, "u"), dbx.NewExp("p.unit_id = u.id")).
			LeftJoin(tableWithAlias(tableNamePhotos, "ph"), dbx.NewExp("p.main_photo_id = ph.id")).
			LeftJoin(tableWithAlias(tableNameProductViewCount, "vc"), dbx.NewExp("p.id = vc.product_id")).
			InnerJoin(tableWithAlias(tableNameGroup, "g"), dbx.NewExp("p.group_id = g.id")).
			InnerJoin(tableWithAlias(tableNameCurrency, "cr"), dbx.NewExp("p.currency_id = cr.id")).
		Where(dbx.NewExp("p.status={:status} and p.group_id={:group_id}", dbx.Params{"status": 1, "group_id": groupId})).
			OrderBy("case when \"vc\".\"count\" is null then 1 else 0 end, \"vc\".\"count\" desc").
			Offset(int64(offset)).
			Limit(int64(limit)).
		All(&rows)

	return rowsToProductEntities(rows), err
}

func (r ProductReadRepository) GetGroupByProductId(ctx context.Context, productId int) (*entity.Group, error) {
	var (
		group Group
	)

	err := r.db.Select("g.id group_id", "g.name group_name", "g.description group_description").
		From(tableWithAlias(tableNameGroup, "g")).
		InnerJoin(tableWithAlias(tableNameProduct, "p"), dbx.NewExp("p.group_id = g.id")).
		Where(dbx.NewExp("p.id={:id}", dbx.Params{"id": productId})).
		One(&group)


	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return &entity.Group{
		ID: group.GroupID,
		Name: group.GroupName,
		Description: group.GroupDescription,
	}, err
}

func (r ProductReadRepository) GetPopularByGroupIdsCount(ctx context.Context, groupIds []int) (int, error) {

	var count int

	ids := make([]interface{}, len(groupIds))

	for k, v := range groupIds {
		ids[k] = v
	}

	err := r.db.Select("COUNT(*)").From(tableWithAlias(tableNameProduct, "p")).
		Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
		Where(dbx.In("p.group_id", ids...)).
		Row(&count)

	return count, err
}

func (r ProductReadRepository) GetPopularByGroupIds(ctx context.Context, groupIds []int, offset, limit int) ([]*entity.Product, error) {

	var rows []Product

	ids := make([]interface{}, len(groupIds))

	for k, v := range groupIds {
		ids[k] = v
	}

	err := r.db.Select(
			"b.id brand_id", "b.name brand_name", "b.slug brand_slug",
			"c.id category_id", "c.name category_name", "c.description category_description", "c.title category_title", "c.slug category_slug",
			"g.id group_id", "g.name group_name", "g.description group_description",
			"cnt.id country_id", "cnt.name country_name",
			"u.id unit_id", "u.name unit_name",
			"ph.id photo_id", "ph.file photo_file", "ph.product_id photo_product_id", "ph.sort photo_sort",
			"cr.id currency_id", "cr.name currency_name", "cr.rate currency_rate", "cr.iso currency_iso",
			"p.*").
		From(tableWithAlias(tableNameProduct, "p")).
			LeftJoin(tableWithAlias(tableNameBrands, "b"), dbx.NewExp("p.brand_id = b.id")).
			LeftJoin(tableWithAlias(tableNameCategories, "c"), dbx.NewExp("p.category_id = c.id")).
			LeftJoin(tableWithAlias(tableNameCountry, "cnt"), dbx.NewExp("p.country_id = cnt.id")).
			LeftJoin(tableWithAlias(tableNameProductUnit, "u"), dbx.NewExp("p.unit_id = u.id")).
			LeftJoin(tableWithAlias(tableNamePhotos, "ph"), dbx.NewExp("p.main_photo_id = ph.id")).
			LeftJoin(tableWithAlias(tableNameProductViewCount, "vc"), dbx.NewExp("p.id = vc.product_id")).
			InnerJoin(tableWithAlias(tableNameGroup, "g"), dbx.NewExp("p.group_id = g.id")).
			InnerJoin(tableWithAlias(tableNameCurrency, "cr"), dbx.NewExp("p.currency_id = cr.id")).
		Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
		Where(dbx.In("p.group_id", ids...)).
			OrderBy("case when \"vc\".\"count\" is null then 1 else 0 end, \"vc\".\"count\" desc").
			Offset(int64(offset)).
			Limit(int64(limit)).
		All(&rows)

	return rowsToProductEntities(rows), err
}

func (r ProductReadRepository) GetGroupsByProductIds(ctx context.Context, productIds []int) ([]*entity.Group, error) {
	var (
		rows []Group
	)

	ids := make([]interface{}, len(productIds))

	for k, v := range productIds {
		ids[k] = v
	}

	err := r.db.Select("g.id group_id", "g.name group_name", "g.description group_description").
		From(tableWithAlias(tableNameGroup, "g")).
		InnerJoin(tableWithAlias(tableNameProduct, "p"), dbx.NewExp("p.group_id = g.id")).
		Where(dbx.In("p.id", ids...)).
		All(&rows)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	groups := make([]*entity.Group, len(rows))

	for k, v := range rows {
		groups[k] = &entity.Group{
			ID:          v.GroupID,
			Name:        v.GroupName,
			Description: v.GroupDescription,
		}
	}

	return groups, err
}

func (r ProductReadRepository) Search(ctx context.Context, input string, size int) ([]*entity.Product, error) {
	var (
		buf bytes.Buffer
		split []string
	)

	name := strings.ToLower(input)
	split = strings.Split(name, " ")

	should := make([]interface{}, len(split) + 1)

	should[0] = map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": name,
			"fields": []string{
				"name",
			},
		},
	}

	for k, v := range split {
		should[k+1] = map[string]interface{}{
			"wildcard": map[string]string{
				"name": "*" + v + "*",
			},
		}
	}

	q := map[string]interface{}{
		"_source": []string{
			"id",
		},
		"size": size,
		"sort": []interface{}{
			map[string]interface{}{
				"_score": []interface{}{
					map[string]string{
						"order": "desc",
					},
				},
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("shop"),
		r.es.Search.WithDocumentType("products"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithPretty(),
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

	var (
		result map[string]interface{}
	)

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return nil, err
	}

	i := 0
	ids := make([]int, len(result["hits"].(map[string]interface{})["hits"].([]interface{})))

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {

		ids[i] = stringIdToInt(fmt.Sprintf("%v", hit.(map[string]interface{})["_id"]))
		i++
	}

	return r.GetByIdsWithSequence(ctx, ids)
}

func (r ProductReadRepository) Exist(ctx context.Context, id int) (bool, error) {
	var count int

	err := r.db.Select("COUNT(*)").From(tableWithAlias(tableNameProduct, "p")).
		Where(dbx.NewExp("p.status={:status}", dbx.Params{"status": 1})).
		Where(dbx.NewExp("p.id={:id}", dbx.Params{"id": id})).
		Row(&count)

	return count > 0, err
}

func (r ProductReadRepository) GetGroupById(ctx context.Context, id int) (*entity.Group, error) {
	var (
		group Group
	)

	err := r.db.Select("g.id group_id", "g.name group_name", "g.description group_description").
		From(tableWithAlias(tableNameGroup, "g")).
		Where(dbx.NewExp("g.id={:id}", dbx.Params{"id": id})).
		One(&group)


	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return &entity.Group{
		ID: group.GroupID,
		Name: group.GroupName,
		Description: group.GroupDescription,
	}, err
}

func (r ProductReadRepository) GetGroupsByIds(ctx context.Context, ids []int) ([]*entity.Group, error) {

	var (
		rows []Group
	)

	groupIds := make([]interface{}, len(ids))

	for k, v := range ids {
		groupIds[k] = v
	}

	err := r.db.Select("g.id group_id", "g.name group_name", "g.description group_description").
		From(tableWithAlias(tableNameGroup, "g")).
		Where(dbx.In("g.id", groupIds...)).
		All(&rows)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	groups := make([]*entity.Group, len(rows))

	for k, v := range rows {
		groups[k] = &entity.Group{
			ID:          v.GroupID,
			Name:        v.GroupName,
			Description: v.GroupDescription,
		}
	}

	return groups, err
}