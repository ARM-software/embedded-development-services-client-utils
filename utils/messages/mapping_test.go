package messages

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
)

func Test_convertRawMessageIntoIMessage(t *testing.T) {
	tests := []struct {
		message any
	}{
		{
			message: *client.NewMessageObject(faker.Sentence()),
		},
		{
			message: *client.NewNotificationMessageObject(faker.Sentence()),
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			require.NotNil(t, test.message)
			_, err := convertRawMessageIntoIMessage(&test.message)
			require.NoError(t, err)
		})
	}
}
