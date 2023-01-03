package cache

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ServiceCache(t *testing.T) {
	cache := NewServiceCache()
	assert.Empty(t, cache.GetKey())
	key := faker.UUIDHyphenated()
	require.NoError(t, cache.SetKey(key))
	assert.Equal(t, key, cache.GetKey())
	require.NoError(t, cache.Invalidate(context.TODO()))
	assert.Empty(t, cache.GetKey())
	require.NoError(t, cache.SetCacheControl(NoStore))
	require.NoError(t, cache.SetKey(key))
	assert.Empty(t, cache.GetKey())
}
