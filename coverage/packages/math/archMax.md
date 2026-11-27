# Independent review of `math.archMax`

**Verdict: FAITHFUL**

Evidence SHA-256: `14f710776ab831c5805fe3cf34fa4f8b9b0d794826e2497fdc6bc1db68ea46ce`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fc6fe4c4a493cf1609993dde4260a3fec159d109eca852f581259a86fd67cfc7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `23d49b3a4ec24ce6308fbab6582179237b305393f0b284b8ed0192119d7aa0b7`.

Pinned Go 1.27.0 tests positive infinity by raw word before processing any NaN. On amd64 the assembly returns its exact canonical NaN and prefers positive zero for a mixed-sign zero pair. On arm64 the hardware instruction gives priority to signaling NaNs and then takes the instruction's first operand among two NaNs of the same class. Go places the second source argument in that operand. The result preserves sign and payload while quieting a selected signaling NaN. The no-assembly entry panics with `not implemented`. [Model](../../../lean/GoLean/Math/ArchDim.lean) · [Comparison](../../../difftest/math/arch-dim/check.go)

The pinned darwin/arm64 overlay matched 784 boundaries and 60,000 samples from seeds 17, 2718 and 65537 per declaration. Independent controls matched 13 native pairs and 39 Lean profile pairs, including swapped NaNs, mixed signed zeros, subnormals and dominant infinities. Canonicalizing ARM NaNs, selecting the first source NaN, and placing NaNs ahead of dominant infinity failed eight, five and two paired fixtures respectively. Native amd64 and no-assembly execution, other assembly architectures, and caller-modified floating-point control remain outside local verification. [Reviewer control](../../../difftest/math/arch-dim/review/main.go)
