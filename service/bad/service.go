package bad

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

type RawItemsSearcher interface {
	SearchByName(name string) ([]*ds.RawItem, error)
}

type Service struct {
	repo             Repo
	itemNamesGetter  ItemNamesGetter
	rawItemsSearcher RawItemsSearcher
}

func NewService(repo Repo, itemNamesGetter ItemNamesGetter, rawItemsSearcher RawItemsSearcher) *Service {
	return &Service{
		repo:             repo,
		itemNamesGetter:  itemNamesGetter,
		rawItemsSearcher: rawItemsSearcher,
	}
}

func (s *Service) Search(id int64, name string) ([]*ds.Item, error) {
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
	ids, names, err := s.searchByName(name)
	if err != nil {
		return nil, err
	}

	validIDs, err := s.repo.GetItemIDs(ids)
	if err != nil {
		return nil, err
	}

	return ds.CreateItems(validIDs, names), nil
}

func (s *Service) searchByName(name string) ([]int64, map[int64]string, error) {
	rawItems, err := s.rawItemsSearcher.SearchByName(name)
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
