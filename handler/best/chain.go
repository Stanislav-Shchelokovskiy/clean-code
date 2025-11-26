package good

import (
	"context"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

type Handler interface {
	Handle(ctx context.Context, variant ds.Variant, diff ds.Diff) (bool, error)
}

type chain struct {
	handlers []Handler
}

func NewChain(itemsUpdater ItemsUpdater, updaters []Updater) *chain {
	handlers := make([]Handler, 0, len(updaters))
	for _, updater := range updaters {
		handlers = append(handlers, NewHandler(itemsUpdater, updater))
	}
	return &chain{handlers: handlers}
}

func (c *chain) Handle(ctx context.Context, e ds.VariantChangeEvent) error {
	for _, variant := range e.Variants {
		if len(variant.ItemIDs) == 0 {
			continue
		}

		for _, diff := range e.Diff {
			for _, h := range c.handlers {
				handled, err := h.Handle(ctx, variant, diff)
				if err != nil {
					return err
				}
				if handled {
					break
				}
			}
		}
	}

	return nil
}
