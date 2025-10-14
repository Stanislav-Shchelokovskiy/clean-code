package good

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

type Searcher interface {
	SearchByName(name string) ([]int64, map[int64]string, error)
}

type Service struct {
	repo            Repo
	itemNamesGetter ItemNamesGetter
	itemsSearcher   Searcher
}

func NewService(repo Repo, itemNamesGetter ItemNamesGetter, rawItemsSearcher RawItemsSearcher) *Service {
	return newService(repo, itemNamesGetter, newRawItemsSearcher(rawItemsSearcher))
}

func newService(repo Repo, itemNamesGetter ItemNamesGetter, searcher Searcher) *Service {
	return &Service{
		repo:            repo,
		itemNamesGetter: itemNamesGetter,
		itemsSearcher:   searcher,
	}
}

type internalService interface {
	SearchByIDs(ids []int64) ([]*ds.Item, error)
	SearchByName(name string) ([]*ds.Item, error)
}

func (s *Service) Search(id int64, name string) ([]*ds.Item, error) {
	return search(s, id, name)
}

func search(s internalService, id int64, name string) ([]*ds.Item, error) {
	if id > 0 {
		return s.SearchByIDs([]int64{id})
	}
	return s.SearchByName(name)
}

func (s *Service) SearchByIDs(ids []int64) ([]*ds.Item, error) {
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

func (s *Service) SearchByName(name string) ([]*ds.Item, error) {
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
