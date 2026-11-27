# Independent review of `math.archExpm1`

**Verdict: FAITHFUL**

Evidence SHA-256: `b59b2d8026c9ce74275e7d039e27e670adeac9f632b890cbd13ba5b9a207817d`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fb87fbbfe8947c10b7cb3d8cd54ed4786ca143073648cc208f79339f2c186ac7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4fd16f76bbd91afe4624d02e50a4a2956f582119afce142adb58b9e5d403f4ad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df1e8bb781cdf3209d26dd3293a2eb952a0857cbc3ef57e30dafec4893a887a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dfc91cf8401dfa86334382300a9e49ac198ade511b85cd859141718186c3e45e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fa71bf08f9d7c155f14eeb3015ac94d3a0ca45c5a69dda4be440ae7178d15b7c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c236bfd996cd02b870c478e1baf5cbfa4bee6f314be382ee08281fd143edec6c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `162c75f787e1c3b9add85929afaf8ff5ba7cb05f7ee8c7e20eb85915090d7b7a`.

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
