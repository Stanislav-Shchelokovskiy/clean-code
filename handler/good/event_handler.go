package good

import (
	"context"
	"encoding/json"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

type Repo interface {
	GetItems(ctx context.Context, ids []int64) ([]*ds.Item, error)
	UpdateItems(ctx context.Context, items []*ds.Item) error
}

type Config interface {
	IgnoreErrors() bool
}

type Updater interface {
	CanUpdate(diff ds.Diff) bool
	Update(ctx context.Context, variant ds.Variant, existingItems []*ds.Item) []*ds.Item
}

type Chain interface {
	Handle(ctx context.Context, e ds.VariantChangeEvent) error
}

type EventHandler struct {
	config Config
	chain  Chain
}

func NewEventHandler(config Config, repo Repo, updaters ...Updater) *EventHandler {
	return newEventHandler(config,
		NewChain(
			NewItemsUpdater(repo),
			updaters,
		),
	)
}

func newEventHandler(config Config, chain Chain) *EventHandler {
	return &EventHandler{
		config: config,
		chain:  chain,
	}
}

func (h *EventHandler) Handle(ctx context.Context, msg ds.Message) (err error) {
	defer func() {
		if err != nil && h.config.IgnoreErrors() {
			err = nil
		}
	}()

	message := ds.VariantChangeEvent{}
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		return nil
	}

	return h.chain.Handle(ctx, message)
}
