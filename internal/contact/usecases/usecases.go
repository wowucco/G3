package usecases

import (
	"context"
	"github.com/wowucco/G3/internal/contact"
	"github.com/wowucco/G3/internal/product"
	"github.com/wowucco/G3/pkg/notification"
)

func NewContactUseCase(n *notification.Service, pr product.Repository) *ContactUserCase {
	return &ContactUserCase{n, pr}
}

type ContactUserCase struct {
	notify *notification.Service
	productRepo product.Repository
}

func (c *ContactUserCase) Recall(ctx context.Context, form contact.IRecallForm) error {

	go c.notify.Recall(form.GetPhone(), form.GetMessage())

	return nil
}

func (c *ContactUserCase) BuyOnClick(ctx context.Context, form contact.IBuyOnClickForm) error {

	p, err := c.productRepo.Get(ctx, form.GetProductId())

	if err != nil {
		return err
	}

	go c.notify.BuyOnClick(form.GetPhone(), *p)

	return nil
}