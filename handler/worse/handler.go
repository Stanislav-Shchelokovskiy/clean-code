package worse

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
	brandAttrID     = 1
	typeAttrID      = 2
)

type Repo interface {
	GetItems(ctx context.Context, ids []int64) ([]*ds.Item, error)
}

type ItemsUpdater interface {
	UpdateItems(ctx context.Context, items []*ds.Item) error
}

type Config interface {
	IgnoreErrors() bool
}

type EventHandler struct {
	config       Config
	repo         Repo
	itemsUpdater ItemsUpdater
}

func NewEventHandler(config Config, repo Repo, itemsUpdater ItemsUpdater) *EventHandler {
	return &EventHandler{
		config:       config,
		repo:         repo,
		itemsUpdater: itemsUpdater,
	}
}

func (h *EventHandler) Handle(ctx context.Context, msg ds.Message) (err error) {
	defer func() {
		if err != nil && h.config.IgnoreErrors() {
			err = nil
		}
	}()

	message := VariantChangeEvent{}
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		return nil
	}
	return h.handle(ctx, message)
}

func (h *EventHandler) handle(ctx context.Context, e VariantChangeEvent) error {
	for _, variant := range e.Variants {
		if len(variant.ItemIDs) == 0 {
			continue
		}

		ids := variant.ItemIDs

		var hasCategoryUpdates bool
		var hasBrandUpdate bool
		var newBrandID int64
		var newTypeID int64

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
		ATTRIBUTES_BRAND_LOOP:
			for _, attr := range variant.Attributes {
				if attr.ID == brandAttrID {
					for _, val := range attr.Values {
						newBrandID = val
						break ATTRIBUTES_BRAND_LOOP
					}
				}
			}

			for _, item := range existingItems {
				if item.BrandID == newBrandID {
					continue
				}

				item.BrandID = newBrandID
				items = append(items, item)
			}

			err := h.itemsUpdater.UpdateItems(ctx, items)
			if err != nil {
				return err
			}
		}

		if hasCategoryUpdates {
			items = items[:0]
			catsByLevels := make(map[int64]int64)
			for _, category := range variant.Categories {
				catsByLevels[category.Level] = category.ID
			}

		ATTRIBUTES_TYPE_LOOP:
			for _, attr := range variant.Attributes {
				if attr.ID == typeAttrID {
					for _, val := range attr.Values {
						newTypeID = val
						break ATTRIBUTES_TYPE_LOOP
					}
				}
			}

			for _, item := range existingItems {
				if item.ParentCategoryID == catsByLevels[1] &&
					item.CategoryID == catsByLevels[2] &&
					item.TypeID == newTypeID {
					continue
				}

				item.ParentCategoryID = catsByLevels[1]
				item.CategoryID = catsByLevels[2]
				item.TypeID = newTypeID
				items = append(items, item)
			}

			if err := h.itemsUpdater.UpdateItems(ctx, items); err != nil {
				return err
			}
		}
	}

	return nil
}
