module github.com/ARM-software/embedded-development-services-client-utils/utils

go 1.27

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
