package uiter

import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	input := slices.Values([]int{0, 1, 2, 3, 4, 5})
	result := Map(Map(input, func(x int) int { return x * x }), strconv.Itoa)
	require.Equal(t, []string{"0", "1", "4", "9", "16", "25"}, Dump(result))
}
