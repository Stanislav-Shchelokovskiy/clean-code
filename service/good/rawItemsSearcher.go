package good

import "github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"

type RawItemsSearcher interface {
	SearchByName(name string) ([]*ds.RawItem, error)
}

type rawItemsSearcher struct {
	itemsSearcher RawItemsSearcher
}

func newRawItemsSearcher(itemsSearcher RawItemsSearcher) *rawItemsSearcher {
	return &rawItemsSearcher{
		itemsSearcher: itemsSearcher,
	}
}

func (s *rawItemsSearcher) SearchByName(name string) ([]int64, map[int64]string, error) {
	rawItems, err := s.itemsSearcher.SearchByName(name)
	if err != nil {
		return nil, nil, err
	}

	ids := make([]int64, 0, len(rawItems))
	names := make(map[int64]string, len(rawItems))
	for _, rawItem := range rawItems {
		ids = append(ids, rawItem.ID)
		names[rawItem.ID] = rawItem.Name
	}

	return ids, names, nil
}
