package contact

import (
	"context"
)

type IContactUseCase interface {
	Recall(ctx context.Context, form IRecallForm) error
	BuyOnClick(ctx context.Context, form IBuyOnClickForm) error
}
