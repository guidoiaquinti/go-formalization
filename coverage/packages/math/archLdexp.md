# Independent review of `math.archLdexp`

**Verdict: FAITHFUL**

Evidence SHA-256: `4bf95a197be8375ae1e9dccad6414c721546709ac8d984af806935b8e6c38541`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c416d0908c98103a700e5b63b6ed73c666303b39d0d5f7e2c89e04018ac790d6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4e9489e91ebc7ac6e00449c81feeb0c6d23bf182bb065202bdd71c5aafcdf261`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b30ffe2278ac4026dbdc9f9f7324bb421bbd5f9e1f6fa319649a5bcceec0dc7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a97c243284343bf96777b962ffbf6b12ff0ddca36217d42304f68398582fd262`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2926d8a7294fa2d1d71366b75a4cb4991e5db016dbe88a3ea6a72c4eaaafa0af`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90376a349711f8d5f2567632e8ac0d13b05e2ca76087412399052185e70716dc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `406431a64e64d447969f2cc1aec0f5856912ef960f048fcbf7ddaba64c1e4554`.

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
