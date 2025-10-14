package bad

import (
	"errors"
	"testing"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/bad/mocks"
	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
	"github.com/stretchr/testify/require"
)

func TestSearchByIDs(t *testing.T) {
	t.Parallel()

	tstErr := errors.New("tstErr")
	ids := []int64{1, 2, 3}
	validIDs := []int64{1, 2}
	names := map[int64]string{1: "name"}
	items := []*ds.Item{{ID: 1, Name: "name"}, {ID: 2}}
	itemsWOnames := []*ds.Item{{ID: 1}, {ID: 2}}

	tests := []struct {
		name      string
		repoErr   error
		getterErr error
		names     map[int64]string
		wantErr   error
		wantItems []*ds.Item
	}{
		{
			name:    "if repo fails, return err",
			repoErr: tstErr,
			wantErr: tstErr,
		},
		{
			name:      "if getter fails, return items without names",
			getterErr: tstErr,
			wantItems: itemsWOnames,
		},
		{
			name:      "if all worked well, return items with found names",
			names:     names,
			wantItems: items,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockRepo(t)
			getter := mocks.NewMockItemNamesGetter(t)

			repo.EXPECT().GetItemIDs(ids).Return(validIDs, test.repoErr).Once()
			if test.repoErr == nil {
				getter.EXPECT().GetNames(validIDs).Return(test.names, test.getterErr).Once()
			}

			s := NewService(repo, getter, nil)
			got, err := s.SearchByIDs(ids)

			require.Equal(t, test.wantErr, err)
			if test.wantErr != nil {
				require.Empty(t, got)
				return
			}
			require.ElementsMatch(t, test.wantItems, got)
		})
	}
}

func TestSearchByNameInner(t *testing.T) {
	t.Parallel()

	tstErr := errors.New("tstErr")
	searchName := "name"
	items := []*ds.RawItem{{ID: 1, Name: "name1"}, {ID: 2, Name: "name2"}}
	wantIDs := []int64{1, 2}
	wantNames := map[int64]string{1: "name1", 2: "name2"}

	tests := []struct {
		name        string
		searcherErr error
		wantErr     error
	}{
		{
			name:        "if searcher fails, return err",
			searcherErr: tstErr,
			wantErr:     tstErr,
		},
		{
			name: "if searcher worked well, return ids and names",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			searcher := mocks.NewMockRawItemsSearcher(t)
			searcher.EXPECT().SearchByName(searchName).Return(items, test.searcherErr)

			s := NewService(nil, nil, searcher)

			gotIDs, gotNames, err := s.searchByName(searchName)

			require.Equal(t, test.wantErr, err)
			if test.wantErr != nil {
				require.Empty(t, gotIDs)
				require.Empty(t, gotNames)
				return
			}
			require.ElementsMatch(t, wantIDs, gotIDs)
			require.Equal(t, wantNames, gotNames)
		})
	}
}

func TestSearchByName(t *testing.T) {
	t.Parallel()

	tstErr := errors.New("tstErr")
	name := "name"
	rawItems := []*ds.RawItem{{ID: 1, Name: "name1"}, {ID: 2, Name: "name2"}}
	ids := []int64{1, 2}
	validIDs := []int64{1}
	items := []*ds.Item{{ID: 1, Name: "name1"}}

	tests := []struct {
		name      string
		repoErr   error
		wantErr   error
		wantItems []*ds.Item
	}{
		{
			name:    "if repo fails, return err",
			repoErr: tstErr,
			wantErr: tstErr,
		},
		{
			name:      "if all worked well, return items",
			wantItems: items,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			searcher := mocks.NewMockRawItemsSearcher(t)
			repo := mocks.NewMockRepo(t)

			searcher.EXPECT().SearchByName(name).Return(rawItems, nil).Once()
			repo.EXPECT().GetItemIDs(ids).Return(validIDs, test.repoErr).Once()

			s := NewService(repo, nil, searcher)
			got, err := s.SearchByName(name)
			require.Equal(t, test.wantErr, err)
			if test.wantErr != nil {
				require.Empty(t, got)
				return
			}
			require.ElementsMatch(t, test.wantItems, got)
		})
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	id := int64(1)
	searchByIDValidIDs := []int64{1}
	searchByIDnames := map[int64]string{1: "name"}
	searchByIDItems := []*ds.Item{{ID: 1, Name: "name"}}

	name := "name"
	rawItems := []*ds.RawItem{{ID: 1, Name: "name1"}, {ID: 2, Name: "name2"}}
	searchByNameIDs := []int64{1, 2}
	searchByNameValidIDs := []int64{1}
	searchByNameItems := []*ds.Item{{ID: 1, Name: "name1"}}

	tests := []struct {
		name       string
		searchID   int64
		names      map[int64]string
		searchName string
		wantItems  []*ds.Item
	}{
		{
			name:      "if id > 0 and all worked well, return items with found names",
			searchID:  id,
			names:     searchByIDnames,
			wantItems: searchByIDItems,
		},
		{
			name:       "if id == 0 and searchName is specified and all worked well, return items",
			searchName: name,
			wantItems:  searchByNameItems,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockRepo(t)
			getter := mocks.NewMockItemNamesGetter(t)
			searcher := mocks.NewMockRawItemsSearcher(t)

			if test.searchID > 0 {
				repo.EXPECT().GetItemIDs([]int64{test.searchID}).Return(searchByIDValidIDs, nil).Once()
				getter.EXPECT().GetNames(searchByIDValidIDs).Return(test.names, nil).Once()
			} else {
				searcher.EXPECT().SearchByName(test.searchName).Return(rawItems, nil).Once()
				repo.EXPECT().GetItemIDs(searchByNameIDs).Return(searchByNameValidIDs, nil).Once()
			}

			s := NewService(repo, getter, searcher)
			got, err := s.Search(test.searchID, test.searchName)

			require.NoError(t, err)
			require.ElementsMatch(t, test.wantItems, got)
		})
	}
}
