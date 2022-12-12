package resourcetests

import (
	"github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/bxcodec/faker/v3"
)

type ResourceTest struct {
	links        *any
	name         string
	title        *string
	resourceType string
}

func (r *ResourceTest) GetLinks() (any, error) {
	if r.links == nil {
		return nil, commonerrors.ErrUndefined
	}
	return *r.links, nil
}

func (r *ResourceTest) GetName() (string, error) {
	return r.name, nil
}

func (r *ResourceTest) GetTitle() (string, error) {
	if r.title == nil {
		return "", commonerrors.ErrUndefined
	}
	return *r.title, nil
}

func (r *ResourceTest) GetType() string {
	return r.resourceType
}

// NewResourceTest generates a fake resource for testing purposes.
func NewResourceTest() (resource.IResource, error) {
	links, err := linkstest.NewFakeLinks()
	if err != nil {
		return nil, err
	}
	title := faker.Name()
	l := any(links)
	return &ResourceTest{
		links:        &l,
		name:         faker.UUIDHyphenated(),
		title:        &title,
		resourceType: faker.Word(),
	}, nil
}
