# Independent review of `math.archLog10`

**Verdict: FAITHFUL**

Evidence SHA-256: `b91db3f54b663391c639ec6df2e890f032d9c92247fdb2e37a6fecca3706b9d7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cd621e69c52bc80368733731a1aa4d63b466cd36b3e985db9b359d4d56ac52b4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9b618bc1fd7bf059f8138ccfd96dbc708c0d6b9d2cb74a69a9d3debad06176f3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1007f46f494193778f996557f2de3a2bfe2d978b11bead32d2ac05681f8bd98d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `45ce7322485b574c94a929a1ec3394888dc2786b292922ad1596afdcbcbb3b78`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bd86017f6dd11a977a3b23cd2750b6454e8e9ecc323a0372feef2157d73bc5e3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a98cdb2bc1b339d92aedc63eef037a811e755adb8ef214ac791c9a8b861a6524`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ea9067f6f97bbc8c8af8da07a34691cdca0794f00a9487838b1e1260e09a164e`.

I checked the pinned Go 1.27.0 declaration, its availability constant and
the non-s390x build constraint. The function body unconditionally panics
with the Go string `not implemented`. The shared independent control
parses the actual pinned source and checks five operand rows plus the
excluded s390x profile for each of the 25 declarations.

The author invoked the exact private symbol with a same-package Go
bridge and checked the panic value and availability flag on darwin/arm64.
Each declaration received operand boundary cases and 60,000 seeded
cases using 17, 2718 and 65537. All 1,503,048 comparisons passed across
the batch. Native non-arm64 ports and the s390x implementations remain
outside the reviewed observations.
