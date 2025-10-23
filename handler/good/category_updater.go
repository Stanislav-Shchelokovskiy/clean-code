package good

import (
	"context"
	"strings"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

const (
	categoryChanged = "Category"
	typeChanged     = "Type"
)

type categoryUpdater struct {
}

// NewBrandUpdater ...
func NewCategoryUpdater() *categoryUpdater {
	return new(categoryUpdater)
}

// CanUpdate ...
func (*categoryUpdater) CanUpdate(diff ds.Diff) bool {
	return strings.HasPrefix(diff.Path, categoryChanged) || strings.HasPrefix(diff.Path, typeChanged)
}

// Update ...
func (u *categoryUpdater) Update(ctx context.Context, variant ds.Variant, existingItems []*ds.Item) []*ds.Item {
	updatedItems := make([]*ds.Item, 0, len(existingItems))
	categories := variant.GetCategories()

	var newTypeID int64
	variant.GetIDs(nil, &newTypeID)

	for _, item := range existingItems {
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
