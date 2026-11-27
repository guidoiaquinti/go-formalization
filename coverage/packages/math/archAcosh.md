# Independent review of `math.archAcosh`

**Verdict: FAITHFUL**

Evidence SHA-256: `4a80d0d6b4ca2674a1ce8b017446080b32a2656c35287d107436d9853ecd2e20`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7c20e1093c40ca27f5eaa340231e96e063b70754d40ebe0549d239efb1951e32`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `20ec6fe5c8e6a40f04a81f4754d9538b6446a2fc500b3d9cd9d5bc851951a8fc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4d00969f60c2f4fd0b7517bab3a44ae163c5d178ccd138a2f82858631b36003`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a412cfafaa56c881cd1a6da8a43cc3e2ba1dd96a1f4b0c403dfaf7af801f28ee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f523d1dde830d8633d6d18b49118f5dfa81515db6d787071e3db76f9a1053279`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `575de3fbb8bf0bbad9d81bfcb4a578803d63f0ed20a20c6468f62ae43614c393`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5c57cc558c56fba1758584ad5cbab3b1e736a38442bfece3fa4ffd81466cdfef`.

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
