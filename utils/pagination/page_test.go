/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package pagination

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
)

func TestToPage(t *testing.T) {
	tests := []struct {
		pageFunc func() (client.IStaticPage, error)
		hasNext  bool
	}{
		{
			pageFunc: func() (client.IStaticPage, error) {
				return nil, nil
			},
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				return client.NewBuildMessageItemWithDefaults(), nil
			},
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				return client.NewBuildMessageItem(*client.NewNullableHalFeedLinks(client.NewHalFeedLinksWithDefaults()), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), []client.MessageObject{*client.NewMessageObject(faker.Sentence()), *client.NewMessageObjectWithDefaults(), *client.NewMessageObject(faker.Word()), *client.NewMessageObjectWithDefaults()}, faker.Word()), nil
			},
			hasNext: false,
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				return client.NewNotificationFeedWithDefaults(), nil
			},
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				return client.NewNotificationFeed(*client.NewNullableHalFeedLinks(client.NewHalFeedLinksWithDefaults()), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), []client.NotificationMessageObject{*client.NewNotificationMessageObject(faker.Sentence()), *client.NewNotificationMessageObjectWithDefaults(), *client.NewNotificationMessageObject(faker.Word()), *client.NewNotificationMessageObjectWithDefaults()}, faker.Word()), nil
			},
			hasNext: false,
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				return client.NewCmsisBuilderCollection(*client.NewNullableHalCollectionLinks(client.NewHalCollectionLinks(*client.NewHalLinkDataWithDefaults())), *client.NewNullablePagingMetadata(client.NewPagingMetadata(0, time.Now(), 50, time.Now(), 45, 100)), faker.Word(), faker.Name()), nil
			},
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				return client.NewCmsisBuilderCollection(*client.NewNullableHalCollectionLinks(&client.HalCollectionLinks{
					Alternate: nil,
					Embedded:  nil,
					First:     nil,
					Item:      []client.HalLinkData{*client.NewHalLinkDataWithDefaults(), *client.NewHalLinkData(faker.URL()), *client.NewHalLinkData(faker.UUIDDigit()), *client.NewHalLinkData(faker.UUIDHyphenated())},
					Last:      nil,
					Next:      client.NewHalLinkDataWithDefaults(),
					Prev:      nil,
					Self:      client.HalLinkData{},
					Simple:    nil,
				}), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), faker.Word(), faker.Name()), nil
			},
			hasNext: true,
		},
		{
			pageFunc: func() (client.IStaticPage, error) {
				builders := client.NewCmsisBuilderCollection(*client.NewNullableHalCollectionLinks(&client.HalCollectionLinks{
					Next: client.NewHalLinkDataWithDefaults(),
					Self: client.HalLinkData{},
				}), *client.NewNullablePagingMetadata(client.NewPagingMetadata(4, time.Now(), 50, time.Now(), 45, 100)), faker.Word(), faker.Name())
				builders.SetEmbedded(client.EmbeddedCmsisBuilderItems{
					Item: []client.CmsisBuilderItem{*client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults()},
				})
				return builders, nil
			},
			hasNext: true,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("page #%v", i), func(t *testing.T) {
			require.NotNil(t, test.pageFunc)
			page, err := test.pageFunc()
			require.NoError(t, err)
			mapped := ToPage(page)
			mappedBack := ToClientPage(mapped)
			if page == nil {
				require.Nil(t, mapped)
				require.Nil(t, mappedBack)
			} else {
				assert.Equal(t, test.hasNext, mapped.HasNext())
				assert.Equal(t, test.hasNext, mappedBack.HasNext())
				count := 0
				it, err := mapped.GetItemIterator()
				require.NoError(t, err)
				itback, err := mappedBack.GetItemIterator()
				require.NoError(t, err)
				for {
					if !it.HasNext() { //nolint:staticcheck
						break
					}
					assert.True(t, itback.HasNext())
					elem, err := it.GetNext()
					assert.NoError(t, err)
					assert.NotNil(t, elem)
					count++
				}
				pageCount, err := mapped.GetItemCount()
				if err != nil {
					assert.Zero(t, pageCount)
				}
				assert.Equal(t, pageCount, int64(count))
				initialPageCount, err := mappedBack.GetItemCount()
				if err != nil {
					assert.Zero(t, initialPageCount)
				}
				assert.Equal(t, pageCount, initialPageCount)

			}
		})
	}
}

func TestUnwrapPage(t *testing.T) {
	var rawPage client.IStaticPage
	assert.Nil(t, rawPage)
	rawPage = client.NewArtefactManagerCollectionWithDefaults()
	require.NotNil(t, rawPage)
	u, ok := rawPage.(*client.ArtefactManagerCollection)
	assert.True(t, ok)
	assert.NotNil(t, u)
	wrappedPage := ToPage(ToClientPage(ToPage(ToClientPage(ToPage(rawPage)))))
	require.NotNil(t, wrappedPage)
	u, ok = wrappedPage.(*client.ArtefactManagerCollection)
	assert.False(t, ok)
	assert.Nil(t, u)
	unwrappedPage := UnwrapPage(wrappedPage)
	require.NotNil(t, unwrappedPage)
	u, ok = unwrappedPage.(*client.ArtefactManagerCollection)
	assert.True(t, ok)
	assert.NotNil(t, u)
}
