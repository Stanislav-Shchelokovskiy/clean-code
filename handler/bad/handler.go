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
	shards       Repo
	config       Config
	itemsUpdater ItemsUpdater
}

func NewEventHandler(config Config, shards Repo, itemsUpdater ItemsUpdater) *EventHandler {
	return &EventHandler{
		shards:       shards,
		config:       config,
		itemsUpdater: itemsUpdater,
	}
}

func (h *EventHandler) Handle(ctx context.Context, msg ds.Message) (err error) {
	defer func() {
		if err != nil && h.config.IgnoreErrors() {
			err = nil
		}
	}()

	message := variantChangeEvent{}
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		return nil
	}
	return h.handle(ctx, message)
}

func (h *EventHandler) handle(ctx context.Context, message variantChangeEvent) error {
	for _, variant := range message.Variants {
		if len(variant.IDs) == 0 {
			continue
		}

		ids := variant.IDs

		var hasCategoryUpdates bool
		var hasBrandUpdate bool
		var newBrandID int64
		var newTypeID int64

		for _, diff := range message.Diff {
			if hasCategoryUpdates && hasBrandUpdate {
				break
			}

			if strings.HasPrefix(diff.Path, categoryChanged) {
				hasCategoryUpdates = true
				continue
			}
			if strings.HasPrefix(diff.Path, brandChanged) {
				hasBrandUpdate = true
				continue
			}
			if strings.HasPrefix(diff.Path, typeChanged) {
				hasCategoryUpdates = true
				continue
			}
		}

		if !hasCategoryUpdates && !hasBrandUpdate {
			continue
		}

		existingItems, err := h.shards.GetItems(ctx, ids)
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

// VariantChangeEvent ...
type variantChangeEvent struct {
	Variants []variant `json:"variants"` // Список товаров
	Diff     []Diff    `json:"diff"`     // Отличия в сравнении
}

type variant struct {
	IDs        []int64     `json:"ids"`        // Список идентификаторов
	Attributes []Attribute `json:"attributes"` // Атрибуты на товаре
	Categories []Category  `json:"categories"` // Категории на товаре
}

type Category struct {
	ID    int64 `json:"id"`    // Идентификатор категории
	Level int64 `json:"level"` // Уровень категории
}
type Attribute struct {
	ID     int64   `json:"id"`     // Идентификатор атрибута
	Values []int64 `json:"values"` // Значения атрибута
}

type Diff struct {
	Path string `json:"path"` // Что изменилось
}
