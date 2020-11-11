package menu

import (
	"context"
	"github.com/wowucco/G3/internal/entity"
)

type ReadRepository interface {

	RootMenuItemWithDepth(ctx context.Context, depth int) (*entity.MenuItem, error)

	MenuItemWithDepthById(ctx context.Context, id, depth int, heedParent bool) (*entity.MenuItem, error)

	Exist(ctx context.Context, id int) (bool, error)
}