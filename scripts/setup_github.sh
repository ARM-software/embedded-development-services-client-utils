#!/bin/bash
##
## Copyright (C) 2020-2026 Arm Limited or its affiliates and Contributors. All rights reserved.
## SPDX-License-Identifier: Apache-2.0
##
go env -w GOPRIVATE=github.com/Arm-Debug
git config --global url."https://${GIT_TOKEN}:x-oauth-basic@github.com/ARM-software".insteadOf "https://github.com/ARM-software"
git config --global url."https://${GIT_TOKEN}:x-oauth-basic@github.com/Arm-Debug".insteadOf "https://github.com/Arm-Debug"
