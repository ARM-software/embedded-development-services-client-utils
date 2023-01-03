/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package job

import (
	"github.com/bxcodec/faker/v3"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests"
)

type MockAsynchronousJob struct {
	resource.IResource
	done    bool
	failure bool
}

func (m *MockAsynchronousJob) FetchType() string {
	return "Mock Asynchronous Job"
}

func (m *MockAsynchronousJob) GetDone() bool {
	return m.done || m.failure
}

func (m *MockAsynchronousJob) GetError() bool {
	return false
}

func (m *MockAsynchronousJob) GetFailure() bool {
	return m.GetDone() && m.failure
}

func (m *MockAsynchronousJob) GetSuccess() bool {
	return m.GetDone() && !m.failure
}

func (m *MockAsynchronousJob) GetStatus() string {
	return faker.Word()
}

func newMockAsynchronousJob(done bool, failure bool) (IAsynchronousJob, error) {
	r, err := resourcetests.NewMockResource()
	if err != nil {
		return nil, err
	}
	return &MockAsynchronousJob{
		IResource: r,
		done:      done,
		failure:   failure,
	}, nil
}

func NewMockUndoneAsynchronousJob() (IAsynchronousJob, error) {
	return newMockAsynchronousJob(false, false)
}

func NewMockSuccessfulAsynchronousJob() (IAsynchronousJob, error) {
	return newMockAsynchronousJob(true, false)
}

func NewMockFailedAsynchronousJob() (IAsynchronousJob, error) {
	return newMockAsynchronousJob(true, true)
}
