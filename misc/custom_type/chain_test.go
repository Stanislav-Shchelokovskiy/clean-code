package custom_type

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	t.Parallel()

	c := &Chain[int]{}
	c.Add([]int{}, []int{1}, []int{}, nil, []int{2, 3}, nil)

	for _, i := range []int{1, 2, 3} {
		require.True(t, c.HasNext())
		require.Equal(t, i, c.Next())
	}
	require.False(t, c.HasNext())
}
