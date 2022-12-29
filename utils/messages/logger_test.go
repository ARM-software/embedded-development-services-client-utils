package messages

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/field"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/logging"
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
)

func TestNewMessageLogger(t *testing.T) {
	defer goleak.VerifyNone(t)
	stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
	require.NoError(t, err)
	logger, err := NewMessageLogger(stdLogger)
	require.NoError(t, err)
	logger.LogEmptyMessageError()
	logger.LogMarshallingError(nil)
	logger.LogMarshallingError(field.ToOptionalAny(*client.NewMessageObject(faker.Sentence())))
	logger.LogMessage(client.NewMessageObject(faker.Sentence()))
	require.NoError(t, logger.Close())
}

func TestLogMessageCollectionCancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
	require.NoError(t, err)
	logger, err := NewMessageLogger(stdLogger)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()
	messages, err := NewMockNotificationFeedPaginator(context.TODO())
	require.NoError(t, err)
	gtx, cancel := context.WithCancel(context.TODO())
	cancel()
	err = logger.LogMessagesCollection(gtx, messages)
	require.Error(t, err)
	assert.True(t, commonerrors.Any(err, commonerrors.ErrCancelled))
	require.NoError(t, logger.Close())
}

func TestLogMessageCollection(t *testing.T) {
	defer goleak.VerifyNone(t)
	stdLogger, err := logging.NewStandardClientLogger(faker.Name(), nil)
	require.NoError(t, err)
	logger, err := NewMessageLogger(stdLogger)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()
	messages, err := NewMockNotificationFeedPaginator(context.TODO())
	require.NoError(t, err)
	require.NoError(t, logger.LogMessagesCollection(context.TODO(), messages))
	require.NoError(t, logger.Close())
}
