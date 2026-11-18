<!--
Copyright (C) 2020-2026 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
# Development and Testing
## Local Development
 
Since this is a Go library, there is no application to run.
 
### Running tests

To run all tests, execute `go test ./...` from each module directory.

To discover the workspace modules from the repository root:
```bash
go list -m -json
```

To refresh workspace and module discovery metadata after layout changes:
```bash
go list -m -json
go work sync
```

### Mocks

We use the built-in mock library and use [mockgen](https://github.com/golang/mock#running-mockgen) to generate the mock implementation of Interfaces for testing.

To generate mocks:
```bash
go generate ./...
```

We make use of `go:generate` comments to instruct which `mockgen` commands should be run.
 
### Coding style
 
`goimports` is enforced via a CI step.
It is up to the individual developer to ensure their change complies with this.

### Static analysis and linting

Static analysis tools and linters are run as part of CI.
They come from [golangci-lint](https://golangci-lint.run/). To run this locally:
```bash
# Run golangci-lint from each Go module directory.
golangci-lint run ./...
``` 

### Precommit

A [pre-commit](https://pre-commit.com/) hook configuration is provided to enforce some tasks on `git commit`.
Run `pre-commit install` to install pre-commit into your git hooks. pre-commit will now run on every commit.

If you want to manually run all pre-commit hooks on a repository without creating hooks, run `pre-commit run --all-files`. 

To run individual hooks use `pre-commit run <hook_id>`.

# Dependency upgrades

For dependency upgrades, dependabot is relied upon and news files are auto-generated in order to document such change.

# Releasing

## Release Types

The CI supports three release flows:

- `development` for snapshot releases
- `release` for stable releases
- `beta` for pre-releases


|   Type      |   Purpose   | Version Number Format | GitHub Release | News Files Deleted |
|-------------|-------------|-----------------------|:--------------:|:------------------:|
| Release     | General Availability | `<minor>.<major>.<patch>`                            | Yes | Yes |
| Beta        | Integration Testing  | `<minor>.<major>.<patch>-beta.<commit number>`       | Yes | No  |
| Development | Development Testing  | `<minor>.<major>.<patch>-dev+<git hash>`             | No  | No  |

> :warning: releases can be made from any branches but
> it is recommended that they are only made from the `main` branch.

### Release workflow

1. Navigate to the [GitHub Actions](https://github.com/ARM-software/embedded-development-services-client-utils/actions/workflows/release.yml) page.
2. Select the **Run Workflow** button and type which kind of release you would like to make (i.e. release, beta or development).

### Version Numbers

The version number will be automatically calculated, based on the news files.

# Detecting secrets

So that no secrets are committed back to the repository, CI runs a combination of two tools:
- `cd-detect-secrets` from [continuous-delivery-scripts](https://github.com/ARMmbed/continuous-delivery-scripts): checks the current repository content for new secret findings against the accepted registry
- [GitLeaks](): scans the git history for usual secrets (for example AWS keys)

To run the repository secret check locally:
```bash
python -m pip install --upgrade continuous-delivery-scripts
cd-detect-secrets
```

Accepted findings are stored in the [baseline](./.secrets.baseline), which `cd-detect-secrets` uses when it runs.

## Baseline & False positives

To flag individual false positives add comment `# pragma: allowlist secret` to the line with the finding.

To update the accepted registry after review, use:
```bash
cd-record-secrets --registry-file .secrets.baseline
```

If you edit `.secrets.baseline` manually on Windows, keep it UTF-8 encoded and convert `\` path separators to `/`.
