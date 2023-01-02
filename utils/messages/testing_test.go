/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
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
		page, err := NewMockNotificationFeedPage(context.TODO(), true, false)
		require.NoError(t, err)
		it, err := page.GetItemIterator()
		require.NoError(t, err)
		for {
			if !it.HasNext() {
				break
			}
			next, err := it.GetNext()
			require.NoError(t, err)
			require.NotNil(t, next)
			_, err = convertRawMessageIntoIMessage(next)
			require.NoError(t, err)
		}
	}
}
