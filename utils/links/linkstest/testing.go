/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package linkstest

import (
	"fmt"

	"github.com/bxcodec/faker/v3"

	"github.com/ARM-software/embedded-development-services-client/client"
)

// NewFakeLink generates a new link for test purposes
func NewFakeLink() (link *any, err error) {
	l := client.NewHalLinkDataWithDefaults()
	err = faker.FakeData(l)
	if err != nil {
		return
	}
	l.SetHref(fmt.Sprintf("link-%x", faker.Word()))
	l2 := any(l)
	link = &l2
	return
}

type FakeLinks struct {
	FakeLink any
}

// NewFakeLinks generates a link structure with fake links.
func NewFakeLinks() (links *FakeLinks, err error) {
	l, err := NewFakeLink()
	if err != nil {
		return
	}
	links = &FakeLinks{
		FakeLink: l,
	}
	return
}
