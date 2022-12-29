package pagination

import (
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
)

func TestToStream(t *testing.T) {
	tests := []struct {
		streamFunc func() (client.IMessageStream, error)
		hasNext    bool
		hasFuture  bool
	}{
		{
			streamFunc: func() (client.IMessageStream, error) {
				return nil, nil
			},
		},
		{
			streamFunc: func() (client.IMessageStream, error) {
				return client.NewBuildMessageItemWithDefaults(), nil
			},
		},
		{
			streamFunc: func() (client.IMessageStream, error) {
				return client.NewBuildMessageItem(*client.NewNullableHalFeedLinks(client.NewHalFeedLinksWithDefaults()), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), []client.MessageObject{*client.NewMessageObject(faker.Sentence()), *client.NewMessageObjectWithDefaults(), *client.NewMessageObject(faker.Word()), *client.NewMessageObjectWithDefaults()}, faker.Word()), nil
			},
			hasNext: false,
		},
		{
			streamFunc: func() (client.IMessageStream, error) {
				return client.NewBuildMessageItem(*client.NewNullableHalFeedLinks(&client.HalFeedLinks{
					First:  nil,
					Future: client.NewHalLinkDataWithDefaults(),
					Last:   nil,
					Next:   client.NewHalLinkDataWithDefaults(),
					Prev:   nil,
					Self:   client.HalLinkData{},
				}), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), []client.MessageObject{*client.NewMessageObject(faker.Sentence()), *client.NewMessageObjectWithDefaults(), *client.NewMessageObject(faker.Word()), *client.NewMessageObjectWithDefaults()}, faker.Word()), nil
			},
			hasNext:   true,
			hasFuture: true,
		},
		{
			streamFunc: func() (client.IMessageStream, error) {
				return client.NewNotificationFeedWithDefaults(), nil
			},
		},
		{
			streamFunc: func() (client.IMessageStream, error) {
				return client.NewNotificationFeed(*client.NewNullableHalFeedLinks(client.NewHalFeedLinksWithDefaults()), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), []client.NotificationMessageObject{*client.NewNotificationMessageObject(faker.Sentence()), *client.NewNotificationMessageObjectWithDefaults(), *client.NewNotificationMessageObject(faker.Word()), *client.NewNotificationMessageObjectWithDefaults()}, faker.Word()), nil
			},
			hasNext: false,
		},
		{
			streamFunc: func() (client.IMessageStream, error) {
				return client.NewNotificationFeed(*client.NewNullableHalFeedLinks(&client.HalFeedLinks{
					First:  nil,
					Future: client.NewHalLinkDataWithDefaults(),
					Last:   nil,
					Next:   client.NewHalLinkDataWithDefaults(),
					Prev:   nil,
					Self:   client.HalLinkData{},
				}), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), []client.NotificationMessageObject{*client.NewNotificationMessageObject(faker.Sentence()), *client.NewNotificationMessageObjectWithDefaults(), *client.NewNotificationMessageObject(faker.Word()), *client.NewNotificationMessageObjectWithDefaults()}, faker.Word()), nil
			},
			hasNext:   true,
			hasFuture: true,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("stream #%v", i), func(t *testing.T) {
			require.NotNil(t, test.streamFunc)
			stream, err := test.streamFunc()
			require.NoError(t, err)
			mapped := ToStream(stream)
			if stream == nil {
				require.Nil(t, mapped)
			} else {
				assert.Equal(t, test.hasNext, mapped.HasNext())
				assert.Equal(t, test.hasFuture, mapped.HasFuture())
				count := 0
				it, err := mapped.GetItemIterator()
				require.NoError(t, err)
				for {
					if !it.HasNext() {
						break
					}
					elem, err := it.GetNext()
					assert.NoError(t, err)
					assert.NotNil(t, elem)
					count += 1
				}
				pageCount, err := mapped.GetItemCount()
				if err != nil {
					assert.Zero(t, pageCount)
				}
				assert.Equal(t, pageCount, int64(count))

			}
		})
	}
}
