module github.com/ARM-software/embedded-development-services-client-utils/utils/cache

go 1.27.1

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
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	golang.org/x/tools v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
