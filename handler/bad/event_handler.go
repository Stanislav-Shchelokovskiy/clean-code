package bad

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

const (
	categoryChanged = "Category"
	brandChanged    = "Brand"
	typeChanged     = "Type"
)

type Repo interface {
	GetItems(ctx context.Context, ids []int64) ([]*ds.Item, error)
	UpdateItems(ctx context.Context, items []*ds.Item) error
}

type Config interface {
	IgnoreErrors() bool
}

type EventHandler struct {
	config Config
	repo   Repo
}

func NewEventHandler(config Config, repo Repo) *EventHandler {
	return &EventHandler{
		config: config,
		repo:   repo,
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
	return h.handle(ctx, message)
}

func (h *EventHandler) handle(ctx context.Context, e ds.VariantChangeEvent) error {
	for _, variant := range e.Variants {
		if len(variant.ItemIDs) == 0 {
			continue
		}

		ids := variant.ItemIDs

		var hasCategoryUpdates bool
		var hasBrandUpdate bool

		for _, diff := range e.Diff {
			if hasCategoryUpdates && hasBrandUpdate {
				break
			}

			if strings.HasPrefix(diff.Path, categoryChanged) || strings.HasPrefix(diff.Path, typeChanged) {
				hasCategoryUpdates = true
				continue
			}

			if strings.HasPrefix(diff.Path, brandChanged) {
				hasBrandUpdate = true
				continue
			}
		}

		if !hasCategoryUpdates && !hasBrandUpdate {
			continue
		}

		existingItems, err := h.repo.GetItems(ctx, ids)
		if err != nil {
			return err
		}

		if len(existingItems) == 0 {
			continue
		}

		items := make([]*ds.Item, 0, len(existingItems))

		if hasBrandUpdate {
			var newBrandID int64
			variant.GetIDs(&newBrandID, nil)
			for _, item := range existingItems {
				if item.SetBrandID(newBrandID) {
					items = append(items, item)
				}
			}

			err := h.repo.UpdateItems(ctx, items)
			if err != nil {
				return err
			}
		}

		if hasCategoryUpdates {
			items = items[:0]
			categories := variant.GetCategories()

			var newTypeID int64
			variant.GetIDs(nil, &newTypeID)

			for _, item := range existingItems {
				if item.SetCategory(
					categories[ds.L1],
					categories[ds.L2],
					newTypeID,
				) {
					items = append(items, item)
				}

				if err := h.repo.UpdateItems(ctx, items); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
