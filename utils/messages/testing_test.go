package messages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newMockHalLink(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, err := newMockHalLink()
		require.NoError(t, err)
	}
}

func Test_NewMockNotificationFeedPage(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, err := NewMockNotificationFeedPage(context.TODO(), true, false)
		require.NoError(t, err)
	}
}
