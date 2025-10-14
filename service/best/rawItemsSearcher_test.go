package best

import (
	"errors"
	"testing"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/bad/mocks"
	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
	"github.com/stretchr/testify/require"
)

func TestRawItemsSearcher(t *testing.T) {
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

			s := newRawItemsSearcher(searcher)

			gotIDs, gotNames, err := s.SearchByName(searchName)

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
