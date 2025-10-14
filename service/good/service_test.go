package good

import (
	"errors"
	"testing"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
	"github.com/Stanislav-Shchelokovskiy/clean-code/service/good/mocks"
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

			s := newService(repo, getter, nil)
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

func TestSearchByName(t *testing.T) {
	t.Parallel()

	tstErr := errors.New("tstErr")
	name := "name"
	ids := []int64{1, 2}
	names := map[int64]string{1: "name1", 2: "name2"}
	validIDs := []int64{1}
	items := []*ds.Item{{ID: 1, Name: "name1"}}

	tests := []struct {
		name        string
		searcherErr error
		repoErr     error
		wantErr     error
		wantItems   []*ds.Item
	}{
		{
			name:        "if searcher fails, return err",
			searcherErr: tstErr,
			wantErr:     tstErr,
		},
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

			searcher := mocks.NewMockSearcher(t)
			repo := mocks.NewMockRepo(t)

			searcher.EXPECT().SearchByName(name).Return(ids, names, test.searcherErr).Once()
			if test.searcherErr == nil {
				repo.EXPECT().GetItemIDs(ids).Return(validIDs, test.repoErr).Once()
			}

			s := newService(repo, nil, searcher)
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

	wantErr := errors.New("tstErr")
	wantItems := []*ds.Item{{ID: 1}}

	tests := []struct {
		name       string
		id         int64
		searchName string
	}{
		{
			name: "if id > 0 , then search by id",
			id:   1,
		},
		{
			name:       "if id == 0, then search by name",
			searchName: "name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewMockInternalService(t)
			if test.id > 0 {
				service.EXPECT().SearchByIDs([]int64{test.id}).Return(wantItems, wantErr)
			} else {
				service.EXPECT().SearchByName(test.searchName).Return(wantItems, wantErr)
			}

			gotItems, err := search(service, test.id, test.searchName)
			require.Equal(t, wantErr, err)
			require.ElementsMatch(t, wantItems, gotItems)
		})
	}
}
