package good

import (
	"context"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

type ItemsUpdater interface {
	Update(ctx context.Context, ids []int64, update func([]*ds.Item) []*ds.Item) error
}

type handler struct {
	itemsUpdater ItemsUpdater
	updater      Updater
}

// NewHandler ...
func NewHandler(itemsUpdater ItemsUpdater, updater Updater) *handler {
	return &handler{
		itemsUpdater: itemsUpdater,
		updater:      updater,
	}
}

// Handle ...
func (h *handler) Handle(ctx context.Context, variant ds.Variant, diff ds.Diff) (bool, error) {
	if !(h.updater.CanUpdate(diff)) {
		return false, nil
	}

	err := h.itemsUpdater.Update(ctx,
		variant.ItemIDs,
		func(existingItems []*ds.Item) []*ds.Item {
			return h.updater.Update(ctx, variant, existingItems)
		},
	)

	return true, err
}
