# Independent review of `math.archMin`

**Verdict: FAITHFUL**

Evidence SHA-256: `8fb4ed105a015885bb2806635217507f219ed7b6dffa390ffb0c8c1f13b34c3b`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `237d19408857aac618e0dec2226d7c1191d013071440e8e8678f4aecd11a9483`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `71173f47e90198ace5f49682bff8d4a7db85c84161aef6019f2e60ae4fdd61a2`.

Pinned Go 1.27.0 tests negative infinity by raw word before processing any NaN. On amd64 the assembly returns its exact canonical NaN and prefers negative zero for a mixed-sign zero pair. On arm64 the hardware instruction gives priority to signaling NaNs and then takes the instruction's first operand among two NaNs of the same class. Go places the second source argument in that operand. The result preserves sign and payload while quieting a selected signaling NaN. The no-assembly entry panics with `not implemented`. [Model](../../../lean/GoLean/Math/ArchDim.lean) · [Comparison](../../../difftest/math/arch-dim/check.go)

The pinned darwin/arm64 overlay matched 784 boundaries and 60,000 samples from seeds 17, 2718 and 65537 per declaration. Independent controls matched 13 native pairs and 39 Lean profile pairs, including swapped NaNs, mixed signed zeros, subnormals and dominant infinities. Canonicalizing ARM NaNs, selecting the first source NaN, and placing NaNs ahead of dominant infinity failed eight, five and two paired fixtures respectively. Native amd64 and no-assembly execution, other assembly architectures, and caller-modified floating-point control remain outside local verification. [Reviewer control](../../../difftest/math/arch-dim/review/main.go)
