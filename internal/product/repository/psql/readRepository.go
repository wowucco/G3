package psql

import (
	"context"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/entity"
)

type ProductReadRepository struct {
	db *dbx.DB
}

func NewProductReadRepository(db *dbx.DB) *ProductReadRepository {
	return &ProductReadRepository{
		db: db,
	}
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