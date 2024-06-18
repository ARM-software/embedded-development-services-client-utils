package job

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/job/jobtest"
	"github.com/ARM-software/embedded-development-services-client/client"
)

func TestAsynchronousJobImplementation(t *testing.T) {
	mock, err := jobtest.NewMockSuccessfulAsynchronousJob()
	require.NoError(t, err)
	tests := []struct {
		impl any
	}{
		{
			impl: client.NewBuildJobItemWithDefaults(),
		},
		{
			impl: client.NewIntellisenseJobItemWithDefaults(),
		},
		{
			impl: client.NewGenericWorkJobItemWithDefaults(),
		},
		{
			impl: mock,
		},
		//{
		//	impl: client.NewVhtRunJobItemWithDefaults(),
		//},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("%T", test.impl), func(t *testing.T) {
			assert.Implements(t, (*IAsynchronousJob)(nil), test.impl)
		})
	}

}
