package product

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
)

type ReadRepository interface {

	GetById(ctx context.Context, id int, with []string) (*entity.Product, error)

	GetByIdsWithSequence(ctx context.Context, ids []int) ([]*entity.Product, error)

	GetSimilar(ctx context.Context, p entity.Product, size int) ([]*entity.Product, error)

	GetPopularCount(ctx context.Context) (int, error)
	GetPopular(ctx context.Context, offset, limit int) ([]*entity.Product, error)

	GetTopSalesCount(ctx context.Context) (int, error)
	GetTopSales(ctx context.Context, offset, limit int) ([]*entity.Product, error)

	GetPopularByGroupIdCount(ctx context.Context, groupId int) (int, error)
	GetPopularByGroupId(ctx context.Context, groupId int, offset, limit int) ([]*entity.Product, error)
	GetGroupByProductId(ctx context.Context, productId int) (*entity.Group, error)

	GetPopularByGroupIdsCount(ctx context.Context, groupIds []int) (int, error)
	GetPopularByGroupIds(ctx context.Context, groupIds []int, offset, limit int) ([]*entity.Product, error)
	GetGroupsByProductIds(ctx context.Context, productIds []int) ([]*entity.Group, error)

	Search(ctx context.Context, input string, size int) ([]*entity.Product, error)
}
