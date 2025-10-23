package good

import (
	"context"
	"strings"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

const brandChanged = "Brand"

type brandUpdater struct{}

// NewBrandUpdater ...
func NewBrandUpdater() *brandUpdater {
	return new(brandUpdater)
}

// CanUpdate ...
func (*brandUpdater) CanUpdate(diff ds.Diff) bool {
	return strings.HasPrefix(diff.Path, brandChanged)
}

// Update ...
func (u *brandUpdater) Update(ctx context.Context, variant ds.Variant, existingItems []*ds.Item) []*ds.Item {
	updatedItems := make([]*ds.Item, 0, len(existingItems))
	var newBrandID int64
	variant.GetIDs(&newBrandID, nil)
	for _, item := range existingItems {
		if item.SetBrandID(newBrandID) {
			updatedItems = append(updatedItems, item)
		}
	}
	return updatedItems
}
