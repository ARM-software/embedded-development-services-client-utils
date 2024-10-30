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

// ToPage converts a page into another
func ToPage(page client.IStaticPage) pagination.IStaticPage {
	return newPageMapper(page)
}

type pageMapper struct {
	page client.IStaticPage
}

func (m *pageMapper) HasNext() bool {
	return m.page.HasNext()
}

func (m *pageMapper) GetItemIterator() (pagination.IIterator, error) {
	iterator, err := m.page.GetItemIterator()
	if err != nil {
		return nil, err
	}
	return ToIterator(iterator), nil
}

func (m *pageMapper) GetItemCount() (int64, error) {
	return m.page.GetItemCount()
}

func newPageMapper(page client.IStaticPage) pagination.IStaticPage {
	if page == nil {
		return nil
	}
	return &pageMapper{
		page: page,
	}
}

// ToClientPage converts a page into another
func ToClientPage(page pagination.IStaticPage) client.IStaticPage {
	return newClientPageMapper(page)
}

type clientPageMapper struct {
	page pagination.IStaticPage
}

func (m *clientPageMapper) HasNext() bool {
	return m.page.HasNext()
}

func (m *clientPageMapper) GetItemIterator() (client.IIterator, error) {
	iterator, err := m.page.GetItemIterator()
	if err != nil {
		return nil, err
	}
	return ToClientIterator(iterator), nil
}

func (m *clientPageMapper) GetItemCount() (int64, error) {
	return m.page.GetItemCount()
}

func newClientPageMapper(page pagination.IStaticPage) client.IStaticPage {
	if page == nil {
		return nil
	}
	return &clientPageMapper{
		page: page,
	}
}

// UnwrapPage tends to unwraps the page mappers to retrieve the raw page object.
func UnwrapPage(page any) any {
	if page == nil {
		return nil
	}
	wrappedPageMapper, ok := page.(*pageMapper)
	if ok {
		return UnwrapPage(wrappedPageMapper.page)
	}
	wrappedClientPageMapper, ok := page.(*clientPageMapper)
	if ok {
		return UnwrapPage(wrappedClientPageMapper.page)
	}
	return UnwrapIterator(page)
}
