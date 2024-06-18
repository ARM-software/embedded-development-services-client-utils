<!--
Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.6.0" (2024-06-18)
=======================

Features
--------

- :sparkles: Improve message logging by checking job affordances (#20240618142953)


Bugfixes
--------

- :bug: [`job`] Wait for job to start before progressing with messages (#20240617160218)


<!--
Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.5.2" (2024-06-17)
=======================

Bugfixes
--------

- Dependency upgrade: upload-artifact-4.3.3 (#20240610175035)
- Dependency upgrade: scorecard-action-2.3.3 (#20240610175041)
- :bug: [`job`] Wait for job to start before progressing with messages (#20240617160218)


<!--
Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.5.1" (2024-06-10)
=======================

Bugfixes
--------

- Dependency upgrade: golang-1.22.0 (#20240214184515)
- Dependency upgrade: fetch-metadata-1.7.0 (#20240321184603)
- Dependency upgrade: golangci-lint-action-4 (#20240422100124)
- Dependency upgrade: golang-1.22.4 (#20240605175520)
- Dependency upgrade: golangci-lint-action-6 (#20240610093118)
- Dependency upgrade: fetch-metadata-2.1.0 (#20240610093119)
- :gear: Upgrade dependencies (#20240610144505)
- :gear: Update deprecated [`faker`](https://github.com/go-faker/faker/) (#20240610144554)


Misc
----

- #202406101045


<!--
Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.5.0" (2024-05-22)
=======================

Features
--------

- :sparkles: Add reusable manager for downloading artefacts (#20240423090834)


Bugfixes
--------

- Dependency upgrade: golang-1.21.4 (#20231108190631)
- Dependency upgrade: setup-python-5 (#20231206181029)
- Dependency upgrade: setup-go-5 (#20231206181033)
- Dependency upgrade: cache-4 (#20240117181416)


Misc
----

- #20240522095850


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.4.1" (2023-10-26)
=======================

Bugfixes
--------

- Dependency upgrade: utils-1.31.1 (#20230426180800)
- Dependency upgrade: atomic-1.11.0 (#20230504180607)
- Dependency upgrade: sync-0.2.0 (#20230504180628)
- Dependency upgrade: fetch-metadata-1.5.0 (#20230522180551)
- Dependency upgrade: goleak-1.3.0 (#20231024171215)
- Dependency upgrade: golang-1.21.3 (#20231026105138)
- Dependency upgrade: checkout-4 (#20231026105824)
- Dependency upgrade: fetch-metadata-1.6.0 (#20231026105857)
- Dependency upgrade: utils-1.51.0 (#20231026122113)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.4.0" (2023-04-21)
=======================

Features
--------

- :sparkle: `[links]` Added a facility for serialising links (#20230421142554)


Bugfixes
--------

- Dependency upgrade: utils-1.31.0 (#20230404180928)
- Dependency upgrade: fetch-metadata-1.4.0 (#20230419180255)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.3.5" (2023-03-20)
=======================

Bugfixes
--------

- Dependency upgrade: checkout-3 (#20230224190552)
- Dependency upgrade: utils-1.28.0 (#20230224191107)
- Dependency upgrade: client-1.21.1 (#20230224191156)
- Dependency upgrade: testify-1.8.2 (#20230228191353)
- Dependency upgrade: setup-go-4 (#20230315190419)
- Dependency upgrade: utils-1.30.0 (#20230320171116)


Deprecations and Removals
-------------------------

- :eraser: Removed the `field` module as it moved to `golang-utils` (#20230320170951)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.3.4" (2023-02-24)
=======================

Bugfixes
--------

- Dependency upgrade: fetch-metadata-1.3.6 (#20230124181152)
- Dependency upgrade: client-1.21.0 (#20230125182059)
- Dependency upgrade: utils-1.26.0 (#20230201182004)
- Dependency upgrade: goleak-1.2.1 (#20230215191333)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.3.3" (2023-01-22)
=======================

Bugfixes
--------

- `[pagination]` extended mappers (#20230120151442)
- `[job]` fixed manager definition (#20230120180053)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.3.2" (2023-01-20)
=======================

Bugfixes
--------

- `[pagination]` extended mappers (#20230120151442)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.3.1" (2023-01-18)
=======================

Bugfixes
--------

- Update `golang-utils` to `1.24.0` (#20230118200409)
- Update `embedded-development-services-client` to `v1.20.0` (#20230118200444)


<!--
Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.3.0" (2023-01-04)
=======================

Features
--------

- :sparkles: `[field]` Added conversion utilities for `any` (#20221229173222)
- :sparkles: `[message]` Added utilities to work with messages (#20221229173247)
- :sparkles: `[pagination]` Added converters for collections and streams (#20221229173318)
- :sparkles: Added `cache` and `store` utilities (#20230103131100)
- :sparkles: Added utilities for dealing with jobs (#20230103145501)


Misc
----

- #20230103131143


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.2.1" (2022-12-23)
=======================

Bugfixes
--------

- Improved error formatting (#20221223152046)
- Use new resource definition (#20221223152507)


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.2.0" (2022-12-22)
=======================

Features
--------

- :sparkle: Utilities to deal with fields (#20221222175115)


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.1.0" (2022-12-14)
=======================

Features
--------

- :sparkles: `[logging]` Added a client helper with various helpers (#20221212173519)


Bugfixes
--------

- Dependency upgrade: utils-1.22.0 (#20221212182051)


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->
"" "1.0.0" (2022-12-09)
=======================

Major changes
-------------

- :sparkles: Public release (#20221209181408)


Bugfixes
--------

- Dependency upgrade: setup-python-4 (#20221209193827)
- Dependency upgrade: utils-1.21.0 (#20221209194039)


Improved Documentation
----------------------

- :book: Added API/Helpers documentation (#20221209191458)


Misc
----

- #202210271204


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Proprietary
-->
"" "0.1.2" (2022-10-24)
=======================

Bugfixes
--------

- Dependency upgrade: testify-1.8.1 (#202210241533)
- Dependency upgrade: client-1.13.0 (#202210241602)


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Proprietary
-->
"" "0.1.1" (2022-10-13)
=======================

Bugfixes
--------

- Change package name (#202210131330)


<!--
Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
SPDX-License-Identifier: Proprietary
-->
"" "0.1.0" (2022-10-13)
=======================

Features
--------

- Add initial library (#202210061512)
