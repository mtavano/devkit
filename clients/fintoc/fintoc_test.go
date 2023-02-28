package fintoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewClient(t *testing.T) {
	cl := NewClient("", "", "", "", nil)
	require.NotNil(t, cl)
}
