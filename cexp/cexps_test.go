package cexp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	require.Equal(t, "true", String(true, "true", "false"))
	require.Equal(t, "false", String(false, "true", "false"))
}

func TestStringShortCircuit(t *testing.T) {
	require.Equal(t, "true", StringShortCircuit(true,
		func() string {
			return "true"
		},
		func() string {
			panic("! not short circuit !")
		},
	))
	require.Equal(t, "false", StringShortCircuit(false,
		func() string {
			panic("! not short circuit !")
		},
		func() string {
			return "false"
		},
	))
}
