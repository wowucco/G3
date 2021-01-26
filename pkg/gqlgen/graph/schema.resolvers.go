package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/pkg/gqlgen/graph/generated"
	"github.com/wowucco/G3/pkg/gqlgen/graph/model"
	"github.com/wowucco/G3/pkg/pagination"
)

func (r *queryResolver) Product(ctx context.Context, input *model.ID) (*model.Product, error) {
	p, e := r.useCase.Get(ctx, input.ID)

	if e != nil {
		return nil, e
	}

	return toProduct(p), nil
}

func (r *queryResolver) Products(ctx context.Context, input *model.Page) (*model.Pages, error) {
	count, _ := r.useCase.Count(ctx)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.useCase.Query(ctx, pages.Offset(), pages.Limit())

	return &model.Pages{
		Page:       pages.Page,
		PerPage:    pages.PerPage,
		PageCount:  pages.PageCount,
		TotalCount: pages.TotalCount,
		Items:      toProducts(ps),
	}, nil
}

func (r *queryResolver) ProductsByIds(ctx context.Context, input *model.Ids) ([]*model.Product, error) {
	ids := make([]int, len(input.Ids))
	for k, v := range input.Ids {
		ids[k] = *v
	}

	ps, err := r.productRead.GetByIdsWithSequence(ctx, ids)

	return toProducts(ps), err
}

func (r *queryResolver) ProductsByGroupID(ctx context.Context, input *model.PageByID) (*model.PagesWithGroup, error) {
	group, err := r.productRead.GetGroupById(ctx, input.ID)

	if err != nil {
		return nil, err
	}

	count, _ := r.productRead.GetPopularByGroupIdCount(ctx, group.ID)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.productRead.GetPopularByGroupId(ctx, group.ID, pages.Offset(), pages.Limit())

	return &model.PagesWithGroup{
		Pages: &model.Pages{
			Page:       pages.Page,
			PerPage:    pages.PerPage,
			PageCount:  pages.PageCount,
			TotalCount: pages.TotalCount,
			Items:      toProducts(ps),
		},
		Group: &model.Group{
			ID:          group.ID,
			Name:        group.Name,
			Description: &group.Description,
		},
	}, nil
}

func (r *queryResolver) ProductsByGroupIds(ctx context.Context, input *model.PageByIds) (*model.PagesWithGroups, error) {
	ids := make([]int, len(input.Ids))

	for k, v := range input.Ids {
		ids[k] = *v
	}

	groups, err := r.productRead.GetGroupsByIds(ctx, ids)

	if err != nil {
		return nil, err
	}

	groupIds := make([]int, len(groups))

	for k, v := range groups {
		groupIds[k] = v.ID
	}

	count, _ := r.productRead.GetPopularByGroupIdsCount(ctx, groupIds)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.productRead.GetPopularByGroupIds(ctx, groupIds, pages.Offset(), pages.Limit())

	groutRows := make([]*model.Group, len(groups))

	for k, v := range groups {
		groutRows[k] = &model.Group{
			ID:          v.ID,
			Name:        v.Name,
			Description: &v.Description,
		}
	}

	return &model.PagesWithGroups{
		Pages: &model.Pages{
			Page:       pages.Page,
			PerPage:    pages.PerPage,
			PageCount:  pages.PageCount,
			TotalCount: pages.TotalCount,
			Items:      toProducts(ps),
		},
		Groups: groutRows,
	}, nil
}

func (r *queryResolver) Popular(ctx context.Context, input *model.Page) (*model.Pages, error) {
	count, _ := r.productRead.GetPopularCount(ctx)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.productRead.GetPopular(ctx, pages.Offset(), pages.Limit())

	return &model.Pages{
		Page:       pages.Page,
		PerPage:    pages.PerPage,
		PageCount:  pages.PageCount,
		TotalCount: pages.TotalCount,
		Items:      toProducts(ps),
	}, nil
}

func (r *queryResolver) Sales(ctx context.Context, input *model.Page) (*model.Pages, error) {
	count, _ := r.productRead.GetTopSalesCount(ctx)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.productRead.GetTopSales(ctx, pages.Offset(), pages.Limit())

	return &model.Pages{
		Page:       pages.Page,
		PerPage:    pages.PerPage,
		PageCount:  pages.PageCount,
		TotalCount: pages.TotalCount,
		Items:      toProducts(ps),
	}, nil
}

func (r *queryResolver) Similar(ctx context.Context, input *model.ID) ([]*model.Product, error) {
	var (
		with []string
	)

	p, err := r.productRead.GetById(ctx, input.ID, with)

	if err != nil {
		return nil, err
	}

	ps, err := r.productRead.GetSimilar(ctx, *p, 10)

	return toProducts(ps), err
}

func (r *queryResolver) Related(ctx context.Context, input *model.ID) ([]*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PopularByProductGroup(ctx context.Context, input *model.PageByID) (*model.PagesWithGroup, error) {
	group, err := r.productRead.GetGroupByProductId(ctx, input.ID)

	if err != nil {
		return nil, err
	}

	count, _ := r.productRead.GetPopularByGroupIdCount(ctx, group.ID)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.productRead.GetPopularByGroupId(ctx, group.ID, pages.Offset(), pages.Limit())

	return &model.PagesWithGroup{
		Pages: &model.Pages{
			Page:       pages.Page,
			PerPage:    pages.PerPage,
			PageCount:  pages.PageCount,
			TotalCount: pages.TotalCount,
			Items:      toProducts(ps),
		},
		Group: &model.Group{
			ID:          group.ID,
			Name:        group.Name,
			Description: &group.Description,
		},
	}, nil
}

func (r *queryResolver) PopularByProductsGroups(ctx context.Context, input *model.PageByIds) (*model.PagesWithGroups, error) {
	ids := make([]int, len(input.Ids))

	for k, v := range input.Ids {
		ids[k] = *v
	}

	groups, err := r.productRead.GetGroupsByProductIds(ctx, ids)

	if err != nil {
		return nil, err
	}

	groupIds := make([]int, len(groups))

	for k, v := range groups {
		groupIds[k] = v.ID
	}

	count, _ := r.productRead.GetPopularByGroupIdsCount(ctx, groupIds)

	pages := pagination.New(input.Page, input.PerPage, count)

	ps, _ := r.productRead.GetPopularByGroupIds(ctx, groupIds, pages.Offset(), pages.Limit())

	groutRows := make([]*model.Group, len(groups))

	for k, v := range groups {
		groutRows[k] = &model.Group{
			ID:          v.ID,
			Name:        v.Name,
			Description: &v.Description,
		}
	}

	return &model.PagesWithGroups{
		Pages: &model.Pages{
			Page:       pages.Page,
			PerPage:    pages.PerPage,
			PageCount:  pages.PageCount,
			TotalCount: pages.TotalCount,
			Items:      toProducts(ps),
		},
		Groups: groutRows,
	}, nil
}

func (r *queryResolver) Search(ctx context.Context, input *model.Text) ([]*model.Product, error) {
	ps, err := r.productRead.Search(ctx, input.Text, 10)

	return toProducts(ps), err
}

func (r *queryResolver) Exist(ctx context.Context, input *model.ID) (*model.ExistProduct, error) {
	exist, err := r.productRead.Exist(ctx, input.ID)

	return &model.ExistProduct{
		Exist: exist,
		ID:    input.ID,
	}, err
}

func (r *queryResolver) TreeMenu(ctx context.Context, input *model.TreeMenu) (*model.TreeMenuItem, error) {
	var (
		depth  int  = 0
		parent bool = false
		id     int  = 0

		e error
		m *entity.MenuItem
	)

	if input != nil {
		if input.ID != nil {
			id = *input.ID
		}

		if input.Depht != nil {
			depth = *input.Depht
		}

		if input.Parent != nil {
			parent = *input.Parent
		}
	}

	if id > 0 {
		m, e = r.menuRead.MenuItemWithDepthById(ctx, id, depth, parent)
	} else {
		m, e = r.menuRead.RootMenuItemWithDepth(ctx, depth)
	}

	if e != nil {
		return nil, e
	}

	return &model.TreeMenuItem{
		ID:          m.Id,
		Name:        m.Name,
		Description: &m.Description,
		Image:       &m.Image,
		Parent:      parentMenuItem(m.Parent),
		Children:    childrenMenuItem(m.Children),
		HasParent:   m.HasParent,
		HasChildren: m.HasChildren,
	}, e
}

func (r *queryResolver) SearchCity(ctx context.Context, input *model.Text) ([]*model.City, error) {
	c, e := r.deliveryRead.SearchCity(ctx, input.Text)

	if e != nil {
		return nil, e
	}

	cities := make([]*model.City, len(c))

	for k, v := range c {
		cities[k] = &model.City{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return cities, nil
}

func (r *queryResolver) CityByID(ctx context.Context, input *model.CityID) (*model.City, error) {
	c, e := r.deliveryRead.GetCityById(ctx, input.ID)

	if e != nil {
		return nil, e
	}

	return &model.City{
		ID:   c.ID,
		Name: c.Name,
	}, nil
}

func (r *queryResolver) DeliveryInfoByCityID(ctx context.Context, input *model.CityID) ([]*model.DeliveryInfo, error) {
	d, e := r.deliveryRead.GetDeliveryInfoByCityId(ctx, input.ID)

	if e != nil {
		return nil, e
	}

	deliveryInfos := make([]*model.DeliveryInfo, len(d))

	for i, value := range d {

		deliveryInfos[i] = &model.DeliveryInfo{
			DeliveryMethod: &model.DeliveryMethod{
				ID:   value.DeliveryMethod.ID,
				Name: value.DeliveryMethod.Name,
				Slug: value.DeliveryMethod.Slug,
			},
			PaymentMethods: deliveryPaymentMethods(value.PaymentMethods),
			Warehouses: deliveryWarehouses(value.Warehouses),
		}
	}

	return deliveryInfos, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.

func deliveryWarehouses(w []entity.Warehouse) []*model.Warehouse {
	
	mw := make([]*model.Warehouse, len(w))
	
	for k, v := range w {
		mw[k] = &model.Warehouse{
			ID:      v.ID,
			Name:    v.Name,
			Address: v.Address,
			Phone:   v.Phone,
		}
	}

	return mw
}

func deliveryPaymentMethods(p []entity.PaymentMethod) []*model.PaymentMethod {

	pm := make([]*model.PaymentMethod, len(p))

	for k, v := range p {
		pm[k] = &model.PaymentMethod{
			ID:   v.ID,
			Name: v.Name,
			Slug: v.Slug,
		}
	}

	return pm
}
func parentMenuItem(parent *entity.ParentMenuItem) *model.TreeParentMenuItem {

	if parent == nil {
		return nil
	}

	row := model.TreeParentMenuItem{
		ID:        parent.Id,
		Name:      parent.Name,
		Image:     &parent.Image,
		HasParent: parent.HasParent,
	}

	if row.HasParent == true {
		row.Parent = parentMenuItem(parent.Parent)
	}

	return &row
}
func childrenMenuItem(children []*entity.ChildrenMenuItem) []*model.TreeChildrenMenuItem {

	ch := make([]*model.TreeChildrenMenuItem, len(children))

	for k, v := range children {
		ch[k] = childMenuItem(v)
	}

	return ch
}
func childMenuItem(child *entity.ChildrenMenuItem) *model.TreeChildrenMenuItem {

	c := &model.TreeChildrenMenuItem{
		ID:          child.Id,
		Name:        child.Name,
		Image:       &child.Image,
		HasChildren: child.HasChildren,
	}

	if c.HasChildren && child.Children != nil {
		c.Children = childrenMenuItem(child.Children)
	}

	return c
}
func toProducts(p []*entity.Product) []*model.Product {
	products := make([]*model.Product, len(p))

	for i, v := range p {
		products[i] = toProduct(v)
	}

	return products
}
func toProduct(p *entity.Product) *model.Product {

	values := make([]*model.CharacteristicValue, len(p.Values))

	for i, val := range p.Values {
		values[i] = &model.CharacteristicValue{
			ID:    val.ID,
			Value: val.Value,
			Characteristic: &model.Characteristic{
				ID:   val.Characteristic.ID,
				Name: val.Characteristic.Name,
				Type: &model.CharacteristicType{
					ID:       val.Characteristic.Type.ID,
					Name:     val.Characteristic.Type.Name,
					IsCustom: val.Characteristic.Type.Custom,
				},
				Unit: &model.Unit{
					ID:   val.Characteristic.Unit.ID,
					Name: val.Characteristic.Unit.Name,
				},
			},
		}
	}

	photos := make([]*model.Photo, len(p.Photos))

	for i, val := range p.Photos {
		photos[i] = &model.Photo{
			ID:     val.ID,
			IsMain: false,
			Sort:   val.Rating,
			Small:  val.GetSmallUrl(p),
			Thumb:  val.GetThumbUrl(p),
		}
	}

	salePrice := p.Price.SaleCentToCurrency()
	return &model.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: &p.Description,
		Code:        p.Code,
		Exist:       p.Exist,
		Status:      p.Status,
		Price: &model.Price{
			Price:            p.Price.CentToCurrency(),
			SalePrice:        &salePrice,
			SaleCount:        &p.Price.SaleCount,
			PriceInCents:     p.Price.Price,
			SalePriceInCents: &p.Price.SalePrice,
			Currency:         p.Price.Currency.Name,
		},
		Brand: &model.Brand{
			ID:   p.Brand.ID,
			Name: p.Brand.Name,
			Slug: p.Brand.Slug,
		},
		Group: &model.Group{
			ID:          p.Group.ID,
			Name:        p.Group.Name,
			Description: &p.Group.Description,
		},
		Category: &model.Category{
			ID:    p.Category.ID,
			Name:  p.Category.Name,
			Title: p.Category.Title,
			Slug:  p.Category.Slug,
		},
		Country: &model.Country{
			ID:   p.Country.ID,
			Name: p.Country.Name,
		},
		Unit: &model.Unit{
			ID:   p.Unit.ID,
			Name: p.Unit.Name,
		},
		MainPhoto: &model.Photo{
			ID:     p.MainPhoto.ID,
			IsMain: p.MainPhoto.IsMain(),
			Sort:   p.MainPhoto.Rating,
			Small:  p.MainPhoto.GetSmallUrl(p),
			Thumb:  p.MainPhoto.GetThumbUrl(p),
		},
		Values: values,
		Photos: photos,
	}
}

type mutationResolver struct{ *Resolver }
