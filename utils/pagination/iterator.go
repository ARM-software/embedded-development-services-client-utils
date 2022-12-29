/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package pagination extends the functionality defined in the `pagination` module present golang-utils
package pagination

import (
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
)

// ToIterator converts an iterator into another
func ToIterator(iterator client.IIterator) pagination.IIterator {
	return newIteratorMapper(iterator)
}

type iteratorMapper struct {
	iterator client.IIterator
}

func (i *iteratorMapper) HasNext() bool {
	return i.iterator.HasNext()
}

func (i *iteratorMapper) GetNext() (*interface{}, error) {
	return i.iterator.GetNext()
}

func newIteratorMapper(iterator client.IIterator) pagination.IIterator {
	if iterator == nil {
		return nil
	}
	return &iteratorMapper{
		iterator: iterator,
	}
}
