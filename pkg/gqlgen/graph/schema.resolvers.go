package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/pkg/gqlgen/graph/generated"
	"github.com/wowucco/G3/pkg/gqlgen/graph/model"
	"github.com/wowucco/G3/pkg/pagination"
)

func (r *queryResolver) Product(ctx context.Context, input *model.GetProductInput) (*model.Product, error) {
	p, _ := r.useCase.Get(ctx, input.ProductID)

	return toProduct(p), nil
}

func (r *queryResolver) Products(ctx context.Context, input *model.GetProductsInput) (*model.Pages, error) {
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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
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
