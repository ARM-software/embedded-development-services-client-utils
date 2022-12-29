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
