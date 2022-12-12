package linkstest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFakeLink(t *testing.T) {
	l, err := NewFakeLink()
	require.NoError(t, err)
	fmt.Println(l)
}
