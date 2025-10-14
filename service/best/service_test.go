package best

import (
	"errors"
	"testing"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/best/mocks"
	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
	"github.com/stretchr/testify/require"
)

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

			byIDSearcher := mocks.NewMockByIDSearcher(t)
			byNameSearcher := mocks.NewMockByNameSearcher(t)
			if test.id > 0 {
				byIDSearcher.EXPECT().SearchByIDs([]int64{test.id}).Return(wantItems, wantErr)
			} else {
				byNameSearcher.EXPECT().SearchByName(test.searchName).Return(wantItems, wantErr)
			}

			s := newService(byIDSearcher, byNameSearcher)
			gotItems, err := s.Search(test.id, test.searchName)

			require.Equal(t, wantErr, err)
			require.ElementsMatch(t, wantItems, gotItems)
		})
	}
}

func TestSearchByIDs(t *testing.T) {
	t.Parallel()

	ids := []int64{1}
	wantErr := errors.New("tstErr")
	wantItems := []*ds.Item{{ID: 1}}

	byIDSearcher := mocks.NewMockByIDSearcher(t)
	byIDSearcher.EXPECT().SearchByIDs(ids).Return(wantItems, wantErr)

	s := newService(byIDSearcher, nil)
	gotItems, err := s.SearchByIDs(ids)

	require.Equal(t, wantErr, err)
	require.ElementsMatch(t, wantItems, gotItems)
}

func TestSearchByName(t *testing.T) {
	t.Parallel()

	searchName := "name"
	wantErr := errors.New("tstErr")
	wantItems := []*ds.Item{{ID: 1}}

	byNameSearcher := mocks.NewMockByNameSearcher(t)
	byNameSearcher.EXPECT().SearchByName(searchName).Return(wantItems, wantErr)

	s := newService(nil, byNameSearcher)
	gotItems, err := s.SearchByName(searchName)

	require.Equal(t, wantErr, err)
	require.ElementsMatch(t, wantItems, gotItems)
}
