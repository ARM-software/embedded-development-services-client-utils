/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package messages

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
)

func TestFormatMessage(t *testing.T) {
	message, err := FormatMessage(nil)
	require.Error(t, err)
	errortest.AssertError(t, err, commonerrors.ErrUndefined)
	assert.Empty(t, message)
	message, err = FormatMessage(client.NewMessageObject(faker.Sentence()))
	require.NoError(t, err)
	assert.NotEmpty(t, message)
	message, err = FormatMessage(client.NewMessageObjectWithDefaults())
	require.NoError(t, err)
	assert.Empty(t, message)

	messageO := client.NewMessageObject(faker.Sentence())
	messageO.SetSource(faker.Name())
	messageO.SetCtime(time.Now())
	messageO.SetSeverity("MAJOR")
	message, err = FormatMessage(messageO)
	require.NoError(t, err)
	assert.NotEmpty(t, message)
	fmt.Println(message)

	message1 := client.NewNotificationMessageObject(faker.Sentence())
	message1.SetSource(faker.Name())
	message1.SetCtime(time.Now())
	message1.SetSeverity("Minor")
	message, err = FormatMessage(message1)
	require.NoError(t, err)
	assert.NotEmpty(t, message)
	fmt.Println(message)

	t.Run("without options", func(t *testing.T) {
		message1 := client.NewNotificationMessageObject(faker.Sentence())
		message1.SetSource(faker.Name())
		message1.SetCtime(time.Now())
		message1.SetSeverity("Minor")
		message, err = FormatMessageWithOptions(message1)
		require.NoError(t, err)
		assert.NotEmpty(t, message)
		fmt.Println(message)
	})
	t.Run("with severity", func(t *testing.T) {
		message1 := client.NewNotificationMessageObject(faker.Sentence())
		message1.SetSource(faker.Name())
		message1.SetCtime(time.Now())
		message1.SetSeverity("Minor")
		message, err = FormatMessageWithOptions(message1, WithSeverity)
		require.NoError(t, err)
		assert.NotEmpty(t, message)
		fmt.Println(message)
	})
}
