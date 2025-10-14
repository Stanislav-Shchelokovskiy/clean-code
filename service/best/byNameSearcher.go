package best

import "github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"

type Searcher interface {
	SearchByName(name string) ([]int64, map[int64]string, error)
}

type byNameSearcher struct {
	repo          Repo
	itemsSearcher Searcher
}

func newByNameSearcher(repo Repo, itemsSearcher Searcher) *byNameSearcher {
	return &byNameSearcher{
		repo:          repo,
		itemsSearcher: itemsSearcher,
	}
}

func (s *byNameSearcher) SearchByName(name string) ([]*ds.Item, error) {
	ids, names, err := s.itemsSearcher.SearchByName(name)
	if err != nil {
		return nil, err
	}

	validIDs, err := s.repo.GetItemIDs(ids)
	if err != nil {
		return nil, err
	}

	return ds.CreateItems(validIDs, names), nil
}
