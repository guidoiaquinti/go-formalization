# Independent review of `math.archErfc`

**Verdict: FAITHFUL**

Evidence SHA-256: `6f8389dd62584ccaab4ec02e9c790f5d2a8876a9bc4bef8bebcdca6f739c4a9d`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d957bcc116c484c1206b0f1601499f2ced5c2cb03d5995cdd9af74a701a7c53f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `42d373e4702e7cd326685d6bd361280b03a6bbe9fa41066722bd812c1a205a6e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cdb744852d2254f7fb4339c39a1bd3de3fad952e96aa5adf0cea8cfaba6ab681`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `514835080ca80aa0d40aa25d6cea45645e4e3a3fcb681763bf21e41c182c3b63`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3c2d549fbbda4357594fc08c2231a2b890f67085d1b1ce8b37e391025e1044c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `104b9c954748c49f5bb073acec713632440c984fec617a2f6a7c82a933fafd2a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fa5e21708fa6e04a94bd6e8a8038aa95b1d11c8ae215c6d974a82c1a932b48b9`.

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
