/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package pagination extends the functionality defined in the `pagination` module present golang-utils
package pagination

import (
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
)

// ToStream converts a message stream into another
func ToStream(stream client.IMessageStream) pagination.IStaticPageStream {
	return newStreamMapper(stream)
}

type streamMapper struct {
	stream client.IMessageStream
}

func (m *streamMapper) HasNext() bool {
	return m.stream.HasNext()
}

func (m *streamMapper) GetItemIterator() (pagination.IIterator, error) {
	iterator, err := m.stream.GetItemIterator()
	if err != nil {
		return nil, err
	}
	return ToIterator(iterator), nil
}

func (m *streamMapper) GetItemCount() (int64, error) {
	return m.stream.GetItemCount()
}

func (m *streamMapper) HasFuture() bool {
	return m.stream.HasFuture()
}

func newStreamMapper(clientStream client.IMessageStream) pagination.IStaticPageStream {
	if clientStream == nil {
		return nil
	}
	return &streamMapper{
		stream: clientStream,
	}
}

// ToClientStream converts a message stream into another
func ToClientStream(stream pagination.IStaticPageStream) client.IMessageStream {
	return newClientStreamMapper(stream)
}

type clientStreamMapper struct {
	stream pagination.IStaticPageStream
}

func (m *clientStreamMapper) HasNext() bool {
	return m.stream.HasNext()
}

func (m *clientStreamMapper) GetItemIterator() (client.IIterator, error) {
	iterator, err := m.stream.GetItemIterator()
	if err != nil {
		return nil, err
	}
	return ToClientIterator(iterator), nil
}

func (m *clientStreamMapper) GetItemCount() (int64, error) {
	return m.stream.GetItemCount()
}

func (m *clientStreamMapper) HasFuture() bool {
	return m.stream.HasFuture()
}

func newClientStreamMapper(stream pagination.IStaticPageStream) client.IMessageStream {
	if stream == nil {
		return nil
	}
	return &clientStreamMapper{
		stream: stream,
	}
}
