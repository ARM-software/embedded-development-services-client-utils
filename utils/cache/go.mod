module github.com/ARM-software/embedded-development-services-client-utils/utils/cache

go 1.27

tool go.uber.org/mock/mockgen

replace github.com/ARM-software/embedded-development-services-client-utils/utils/api => ../api

replace github.com/ARM-software/embedded-development-services-client-utils/utils/artefacts => ../artefacts

replace github.com/ARM-software/embedded-development-services-client-utils/utils/client => ../client

replace github.com/ARM-software/embedded-development-services-client-utils/utils/errors => ../errors

replace github.com/ARM-software/embedded-development-services-client-utils/utils/job => ../job

replace github.com/ARM-software/embedded-development-services-client-utils/utils/job/jobtest => ../job/jobtest

replace github.com/ARM-software/embedded-development-services-client-utils/utils/links => ../links

replace github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest => ../links/linkstest

replace github.com/ARM-software/embedded-development-services-client-utils/utils/logging => ../logging

replace github.com/ARM-software/embedded-development-services-client-utils/utils/messages => ../messages

replace github.com/ARM-software/embedded-development-services-client-utils/utils/mocks => ../mocks

replace github.com/ARM-software/embedded-development-services-client-utils/utils/pagination => ../pagination

replace github.com/ARM-software/embedded-development-services-client-utils/utils/resource => ../resource

replace github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests => ../resource/resourcetests

replace github.com/ARM-software/embedded-development-services-client-utils/utils/store => ../store

require (
	github.com/go-faker/faker/v4 v4.12.0
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.51.0 // indirect
)
