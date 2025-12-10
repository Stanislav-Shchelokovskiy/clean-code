package not_good_but_viable

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
	return handle(ctx, h, h.repo, e)
}

type innerEvenHandler interface {
	CanHandle(diffs []ds.Diff) (bool, bool)
	UpdateBrands(variant ds.Variant, items []*ds.Item) []*ds.Item
	UpdateCategories(variant ds.Variant, items []*ds.Item) []*ds.Item
}

func handle(ctx context.Context, h innerEvenHandler, r Repo, e ds.VariantChangeEvent) error {
	for _, variant := range e.Variants {
		if len(variant.ItemIDs) == 0 {
			continue
		}

		ids := variant.ItemIDs

		hasCategoryUpdates, hasBrandUpdate := h.CanHandle(e.Diff)
		if !hasCategoryUpdates && !hasBrandUpdate {
			continue
		}

		existingItems, err := r.GetItems(ctx, ids)
		if err != nil {
			return err
		}

		if len(existingItems) == 0 {
			continue
		}

		var updatedItems []*ds.Item
		if hasBrandUpdate {
			updatedItems = h.UpdateBrands(variant, existingItems)
		}

		if hasCategoryUpdates {
			updatedItems = h.UpdateCategories(variant, existingItems)
		}

		if len(updatedItems) > 0 {
			if err := r.UpdateItems(ctx, updatedItems); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *EventHandler) CanHandle(diffs []ds.Diff) (bool, bool) {
	var hasCategoryUpdates bool
	var hasBrandUpdate bool
	for _, diff := range diffs {
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

	return hasCategoryUpdates, hasBrandUpdate
}

func (h *EventHandler) UpdateBrands(variant ds.Variant, items []*ds.Item) []*ds.Item {
	updatedItems := make([]*ds.Item, 0, len(items))
	var newBrandID int64
	variant.GetIDs(&newBrandID, nil)
	for _, item := range items {
		if item.SetBrandID(newBrandID) {
			updatedItems = append(updatedItems, item)
		}
	}
	return updatedItems
}

func (h *EventHandler) UpdateCategories(variant ds.Variant, items []*ds.Item) []*ds.Item {
	updatedItems := make([]*ds.Item, 0, len(items))
	categories := variant.GetCategories()

	var newTypeID int64
	variant.GetIDs(nil, &newTypeID)

	for _, item := range items {
		if item.SetCategory(
			categories[ds.L1],
			categories[ds.L2],
			newTypeID,
		) {
			updatedItems = append(updatedItems, item)
		}
	}
	return updatedItems
}
