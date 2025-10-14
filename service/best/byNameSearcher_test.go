package best

import (
	"errors"
	"testing"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/best/mocks"
	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
	"github.com/stretchr/testify/require"
)

func TestByNameSearcher(t *testing.T) {
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

			s := newByNameSearcher(repo, searcher)
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
