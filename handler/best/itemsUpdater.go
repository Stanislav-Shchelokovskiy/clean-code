package good

import (
	"context"

	"github.com/Stanislav-Shchelokovskiy/clean-code/handler/ds"
)

type itemsUpdater struct {
	repo Repo
}

func NewItemsUpdater(repo Repo) *itemsUpdater {
	return &itemsUpdater{repo: repo}
}

func (u *itemsUpdater) Update(ctx context.Context, skus []int64, update func([]*ds.Item) []*ds.Item) error {
	if len(skus) == 0 {
		return nil
	}

	existingItems, err := u.repo.GetItems(ctx, skus)
	if err != nil {
		return err
	}

	if len(existingItems) == 0 {
		return nil
	}

	updatedItems := update(existingItems)

	if len(updatedItems) == 0 {
		return nil
	}

	return u.repo.UpdateItems(ctx, updatedItems)
}
