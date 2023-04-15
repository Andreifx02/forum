package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSubString(t *testing.T) {
	require.Equal(t, 1, FindSubString("abcd", "bb", 1))	
	require.Equal(t, 2, FindSubString("abcd", "c", 0))
	require.Equal(t, -1, FindSubString("abcd", "dcba", 2))
}