/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package pagination

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
)

func TestToIterator(t *testing.T) {
	tests := []struct {
		iteratorFunc func() (client.IIterator, error)
	}{
		{
			iteratorFunc: func() (client.IIterator, error) {
				return nil, nil
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewHalLinkDataIterator(nil)
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewHalLinkDataIterator([]client.HalLinkData{})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewHalLinkDataIterator([]client.HalLinkData{*client.NewHalLinkDataWithDefaults(), *client.NewHalLinkDataWithDefaults(), *client.NewHalLinkDataWithDefaults(), *client.NewHalLinkDataWithDefaults()})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewBuildersIterator(nil)
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewBuildersIterator([]client.CmsisBuilderItem{})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewBuildersIterator([]client.CmsisBuilderItem{*client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults(), *client.NewCmsisBuilderItemWithDefaults()})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewMessagesIterator(nil)
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewMessagesIterator([]client.MessageObject{})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewMessagesIterator([]client.MessageObject{*client.NewMessageObject(faker.Sentence()), *client.NewMessageObjectWithDefaults(), *client.NewMessageObject(faker.Word()), *client.NewMessageObjectWithDefaults(), *client.NewMessageObject(faker.Name())})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewNotificationMessagesIterator(nil)
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewNotificationMessagesIterator([]client.NotificationMessageObject{})
			},
		},
		{
			iteratorFunc: func() (client.IIterator, error) {
				return client.NewNotificationMessagesIterator([]client.NotificationMessageObject{*client.NewNotificationMessageObject(faker.Sentence()), *client.NewNotificationMessageObject(faker.Word()), *client.NewNotificationMessageObject(faker.Name()), *client.NewNotificationMessageObjectWithDefaults()})
			},
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("iterator #%v", i), func(t *testing.T) {
			require.NotNil(t, test.iteratorFunc)
			iterator, err := test.iteratorFunc()
			require.NoError(t, err)
			mapped := ToIterator(iterator)
			mappedBack := ToClientIterator(mapped)
			if iterator == nil {
				require.Nil(t, mapped)
				require.Nil(t, mappedBack)
			} else {
				for {
					if mapped.HasNext() {
						assert.True(t, mappedBack.HasNext())
						elem, err := mapped.GetNext()
						require.NoError(t, err)
						require.NotNil(t, elem)
					} else {
						elem, err := mapped.GetNext()
						require.Error(t, err)
						require.Empty(t, elem)
						assert.False(t, mappedBack.HasNext())
						elem, err = mappedBack.GetNext()
						require.Error(t, err)
						require.Empty(t, elem)
						break
					}
				}
			}
		})
	}
}
