# Independent review of `math.archMod`

**Verdict: FAITHFUL**

Evidence SHA-256: `37597b8354cc8f8b8067be361959d6dcfa6be9a81b6001241ac7cf1714f386a3`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9cbd3d84a571755f7b526dd4cf46a8989603cd06c4343161cf9f6479f6eac59d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3cb826b1e7efa5b9ace883fabd3670346593d9feaef1498c8b38d963a4270b4d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b9baec9040ef33534854d8a5996dc3246899a2db36d4ca903f976aedce700f8c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b8a64552e2c7edf923802e55ac01572ce7f75ef3e4206478d15c3f7275659df`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ac24cd5c4197f2bae818a3ffaddc92433ce2b458ad662c8bd835ea6999f3fb4b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b21cf615cc37f987700c58f2f81f838b76d66e91df700fe6539a97f1804f66a7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `68be12ee674e967e02acfdd22ddb3dae6adf55b14e30ab3a8cde19a53ccdb20f`.

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
