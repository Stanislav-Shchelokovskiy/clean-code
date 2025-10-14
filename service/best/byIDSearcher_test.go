package best

import (
	"errors"
	"testing"

	"github.com/Stanislav-Shchelokovskiy/clean-code/service/best/mocks"
	"github.com/Stanislav-Shchelokovskiy/clean-code/service/ds"
	"github.com/stretchr/testify/require"
)

func TestByIDSearcher(t *testing.T) {
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

			s := newByIDSearcher(repo, getter)
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
