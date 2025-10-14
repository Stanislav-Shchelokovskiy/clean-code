package best

import (
	"fmt"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
)

type Repo interface {
	GetItemIDs(ids []int64) ([]int64, error)
}

type ItemNamesGetter interface {
	GetNames(ids []int64) (map[int64]string, error)
}

type byIDSearcher struct {
	repo            Repo
	itemNamesGetter ItemNamesGetter
}

func newByIDSearcher(repo Repo, itemNamesGetter ItemNamesGetter) *byIDSearcher {
	return &byIDSearcher{
		repo:            repo,
		itemNamesGetter: itemNamesGetter,
	}
}

func (s *byIDSearcher) SearchByIDs(ids []int64) ([]*ds.Item, error) {
	validIDs, err := s.repo.GetItemIDs(ids)
	if err != nil {
		return nil, err
	}

	names, err := s.itemNamesGetter.GetNames(validIDs)
	if err != nil {
		fmt.Println("getter failed", err)
	}

	return ds.CreateItems(validIDs, names), nil
}
