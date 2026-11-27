# Independent review of `math.archRemainder`

**Verdict: FAITHFUL**

Evidence SHA-256: `6cb31bb044b47376c6f21eeb4f1466edaf3abee71dcdc3354e7206ed731f1f81`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `84c21f4010f62311cdceae5a0c8562f2842fe82c480f5a2c7a3102a410f0e8bb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b079f344120ab9fc328a8cabb265e9edfc1e040cc53de0e2627171ad473e96b7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cb4e76d0bb04ebb3d63cccdef53861694754202d40744bdf25b9b85ed6f79b1e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f8f66653dfe85b3f92d32366fba1d7689806b8b226d012939f5c867c461b85c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dc4e41e0decbb6676b24e1093c238d6b45048fddf0e48552e33a2fe90f79f048`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b905a933313bdd93111401285c8c46849672b52d363ac76e887d6fb788da69c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `57aa59f40935efc15cf83df22e2267e7c0a041b3bc11b746ebf995a247421b2f`.

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
