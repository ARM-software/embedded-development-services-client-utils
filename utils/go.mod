module github.com/ARM-software/embedded-development-services-client-utils/utils

go 1.27.1

require (
	github.com/ARM-software/embedded-development-services-client-utils/utils/api v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/artefacts v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/cache v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/client v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/errors v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/job v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/job/jobtest v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/links v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/logging v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/messages v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/mocks v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/pagination v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/resource v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests v0.0.0
	github.com/ARM-software/embedded-development-services-client-utils/utils/store v0.0.0
)

replace (
	github.com/ARM-software/embedded-development-services-client-utils/utils/api => ./api
	github.com/ARM-software/embedded-development-services-client-utils/utils/artefacts => ./artefacts
	github.com/ARM-software/embedded-development-services-client-utils/utils/cache => ./cache
	github.com/ARM-software/embedded-development-services-client-utils/utils/client => ./client
	github.com/ARM-software/embedded-development-services-client-utils/utils/errors => ./errors
	github.com/ARM-software/embedded-development-services-client-utils/utils/job => ./job
	github.com/ARM-software/embedded-development-services-client-utils/utils/job/jobtest => ./job/jobtest
	github.com/ARM-software/embedded-development-services-client-utils/utils/links => ./links
	github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest => ./links/linkstest
	github.com/ARM-software/embedded-development-services-client-utils/utils/logging => ./logging
	github.com/ARM-software/embedded-development-services-client-utils/utils/messages => ./messages
	github.com/ARM-software/embedded-development-services-client-utils/utils/mocks => ./mocks
	github.com/ARM-software/embedded-development-services-client-utils/utils/pagination => ./pagination
	github.com/ARM-software/embedded-development-services-client-utils/utils/resource => ./resource
	github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests => ./resource/resourcetests
	github.com/ARM-software/embedded-development-services-client-utils/utils/store => ./store
)
