package usecases

import (
	"context"
	"github.com/wowucco/G3/internal/contact"
	"github.com/wowucco/G3/pkg/notification"
)

func NewContactUseCase(n *notification.Service) *ContactUserCase {
	return &ContactUserCase{n}
}

type ContactUserCase struct {
	notify *notification.Service
}

func (c *ContactUserCase) Recall(ctx context.Context, form contact.IRecallForm) error {

	go c.notify.Recall(form.GetPhone(), form.GetMessage())

	return nil
}