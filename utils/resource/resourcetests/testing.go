/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package resourcetests

import (
	"github.com/bxcodec/faker/v3"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
)

type MockResource struct {
	links        *any
	name         string
	title        *string
	resourceType string
}

func (r *MockResource) FetchLinks() (any, error) {
	if r.links == nil {
		return nil, commonerrors.ErrUndefined
	}
	return *r.links, nil
}

func (r *MockResource) FetchName() (string, error) {
	return r.name, nil
}

func (r *MockResource) FetchTitle() (string, error) {
	if r.title == nil {
		return "", commonerrors.ErrUndefined
	}
	return *r.title, nil
}

func (r *MockResource) FetchType() string {
	return r.resourceType
}

// NewMockResource generates a fake resource for testing purposes.
func NewMockResource() (resource.IResource, error) {
	links, err := linkstest.NewFakeLinks()
	if err != nil {
		return nil, err
	}
	title := faker.Name()
	l := any(links)
	return &MockResource{
		links:        &l,
		name:         faker.UUIDHyphenated(),
		title:        &title,
		resourceType: faker.Word(),
	}, nil
}
