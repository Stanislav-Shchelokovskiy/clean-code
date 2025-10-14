package best

import "github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"

type ByIDSearcher interface {
	SearchByIDs(ids []int64) ([]*ds.Item, error)
}

type ByNameSearcher interface {
	SearchByName(name string) ([]*ds.Item, error)
}

type Service struct {
	byIDSearcher   ByIDSearcher
	byNameSearcher ByNameSearcher
}

func NewService(repo Repo, itemNamesGetter ItemNamesGetter, rawItemsSearcher RawItemsSearcher) *Service {
	return newService(
		newByIDSearcher(repo, itemNamesGetter),
		newByNameSearcher(repo, newRawItemsSearcher(rawItemsSearcher)),
	)
}

func newService(byIDSearcher ByIDSearcher, byNameSearcher ByNameSearcher) *Service {
	return &Service{
		byIDSearcher:   byIDSearcher,
		byNameSearcher: byNameSearcher,
	}
}

func (s *Service) Search(id int64, name string) ([]*ds.Item, error) {
	if id > 0 {
		return s.byIDSearcher.SearchByIDs([]int64{id})
	}
	return s.byNameSearcher.SearchByName(name)
}

func (s *Service) SearchByIDs(ids []int64) ([]*ds.Item, error) {
	return s.byIDSearcher.SearchByIDs(ids)
}

func (s *Service) SearchByName(name string) ([]*ds.Item, error) {
	return s.byNameSearcher.SearchByName(name)
}
