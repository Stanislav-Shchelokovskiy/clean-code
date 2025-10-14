package ds

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateItems(t *testing.T) {
	t.Parallel()

	ids := []int64{1, 2}
	names := map[int64]string{1: "name1"}
	wantItems := []*Item{{ID: 1, Name: "name1"}, {ID: 2}}

	got := CreateItems(ids, names)
	require.ElementsMatch(t, wantItems, got)
}
