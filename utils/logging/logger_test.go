package logging

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/logs/logstest"
)

func TestClientLogger(t *testing.T) {
	logger, err := NewClientLogger("test client Logger")
	require.NoError(t, err)

	defer func() { _ = logger.Close() }()

	require.NoError(t, logger.AppendLogger(logstest.NewNullTestLogger(), logstest.NewTestLogger(t), logstest.NewStdTestLogger()))
	require.NoError(t, logger.Check())
	sl, err := logs.NewStdLogger(faker.Name())
	require.NoError(t, err)
	require.NoError(t, logger.Append(sl))

	require.NoError(t, logger.SetLogSource(faker.Name()))
	logger.LogRawError(commonerrors.ErrUndefined)
	logger.LogErrorMessage("%v ....", faker.Sentence())
	logger.LogErrorAndMessage(commonerrors.ErrUnexpected, "some error %v (%v): %v", "no idea", 123, faker.Sentence())
	logger.LogInfo("information %v", faker.Sentence())
	resource, err := resourcetests.NewResourceTest()
	require.NoError(t, err)
	logger.LogResource(resource)

	require.NoError(t, logger.Close())
}
