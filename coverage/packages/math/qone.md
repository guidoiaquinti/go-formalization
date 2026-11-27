# Independent review of `math.qone`

**Verdict: FAITHFUL**

Evidence SHA-256: `2ea74cfc31e56c59c5c0e8aacde4c4f59afce5adce66d1748ed59af13275d35a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2399159e9a906337941c1872f26740c57616b01faa5bd43620bcad034bcfe2d8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `85e7ad81b6a2f2ea57441204fc93d85110918b7072c685ab9bf147af0bb14151`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `60954d82c901c814b49564380c6a4c3d8be0977c8f2079ffb3203ca70f892ebb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6f28f2443c6b2d2881ddf1d551e6ee5169fef2512da5c7e11fa34ad1ff6c0aa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `79e7161dfa0127f41ffc5331a16fb736932ed779848ed310c0bc99f19ef92c2a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8657d3ce2b3519d5ff0b5e543f2c2b188f655cd741aa71b04ca18a1703ff7296`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4d6ea62cbc7bb28c075e75a7b8fa2d36aa1950ff1bac3cd37e8646e28703d28a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1a8830b4c13cce0bd3559cbbab8d9162764a8c67b731e52a7e8faf8de297d88d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fc422c27fdaec49c8aa7af1b40d21b7289b92bb54334343cb0b8b84921b9a460`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `27b87c7c3f59ef0eb34feeae8ffdbe8eaf56009c002095beff72e4b1fdaf4f9f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `986f7d8b91124c7ec3264fe5e6cbcebdbb48172f9fb0822f2718ac24952a5df8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bdf5a5e69a8f3d92799fe3a1b3ee0df8da9401c68f526aba3fc0b9d87b0a13a5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4c2f9b0daec5b578ca1bde5d5d5b93836de51ee4974384f449f808f43939a904`.

I checked the pinned Go 1.27.0 private branch selection, coefficient
bands, reciprocal square, denominator leading term and final division.
Arguments below two and NaN reach the nil coefficient table. Positive
infinity returns positive zero.

The independent control evaluated the pinned decimal tables using explicit
rounding and native FMA. It passed 44 exact darwin/arm64 cases, including
eight panics, and 132 Lean observations across three profiles. Omitting the
denominator contraction, omitting Horner contraction, swapping coefficient
bands and reversing the infinity-result zero sign changed 8, 8, 4 and 1
observations respectively. Native amd64 execution remains unverified.

The author comparison passed 66 boundaries and 60,000 seeded observations
using 17, 2718 and 65537, including 10,475 expected panics. The independent
numerical reference differed by at most one binary64 step.
