module github.com/ARM-software/embedded-development-services-client-utils/utils/api

go 1.27

replace github.com/ARM-software/embedded-development-services-client-utils/utils/errors => ../errors

replace github.com/ARM-software/embedded-development-services-client-utils/utils/cache => ../cache

replace github.com/ARM-software/embedded-development-services-client-utils/utils/links => ../links

replace github.com/ARM-software/embedded-development-services-client-utils/utils/artefacts => ../artefacts

replace github.com/ARM-software/embedded-development-services-client-utils/utils/client => ../client

replace github.com/ARM-software/embedded-development-services-client-utils/utils/job => ../job

replace github.com/ARM-software/embedded-development-services-client-utils/utils/job/jobtest => ../job/jobtest

replace github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest => ../links/linkstest

replace github.com/ARM-software/embedded-development-services-client-utils/utils/logging => ../logging

replace github.com/ARM-software/embedded-development-services-client-utils/utils/messages => ../messages

replace github.com/ARM-software/embedded-development-services-client-utils/utils/mocks => ../mocks

replace github.com/ARM-software/embedded-development-services-client-utils/utils/pagination => ../pagination

replace github.com/ARM-software/embedded-development-services-client-utils/utils/resource => ../resource

replace github.com/ARM-software/embedded-development-services-client-utils/utils/resource/resourcetests => ../resource/resourcetests

replace github.com/ARM-software/embedded-development-services-client-utils/utils/store => ../store

require (
	github.com/ARM-software/embedded-development-services-client-utils/utils/errors v0.0.0-00010101000000-000000000000
	github.com/ARM-software/embedded-development-services-client/client v1.109.1
	github.com/ARM-software/golang-utils/utils/commonerrors v1.172.0
	github.com/ARM-software/golang-utils/utils/commonerrors/errortest v1.172.0
	github.com/ARM-software/golang-utils/utils/field v1.172.0
	github.com/ARM-software/golang-utils/utils/http v1.172.0
	github.com/go-faker/faker/v4 v4.12.0
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/ARM-software/golang-utils/utils/collection v1.172.0 // indirect
	github.com/ARM-software/golang-utils/utils/parallelisation v1.172.0 // indirect
	github.com/ARM-software/golang-utils/utils/ptr v1.172.0 // indirect
	github.com/ARM-software/golang-utils/utils/reflection v1.172.0 // indirect
	github.com/ARM-software/golang-utils/utils/safecast v1.172.0 // indirect
	github.com/ARM-software/golang-utils/utils/safeio v1.172.0 // indirect
	github.com/ARM-software/golang-utils/utils/value v1.172.0 // indirect
	github.com/deckarep/golang-set/v2 v2.9.0 // indirect
	github.com/dolmen-go/contextio v1.0.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/petermattis/goid v0.0.0-20250813065127-a731cc31b4fe // indirect
	github.com/sasha-s/go-deadlock v0.3.9 // indirect
	go.mongodb.org/mongo-driver v1.17.7 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)
