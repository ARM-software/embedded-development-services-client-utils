package resourcetests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResourceTest(t *testing.T) {
	l, err := NewResourceTest()
	require.NoError(t, err)
	fmt.Println(l)
}
