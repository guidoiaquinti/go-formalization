# Independent review of `math.qzero`

**Verdict: FAITHFUL**

Evidence SHA-256: `4ff1314929a7b6b1a1a7a35d7a70aaa8a1051099cc7de079981d05baa8026c34`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9e41b753c527f1b4274bc37838a933f74320585a58d303feb75600d23bd5d942`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `44a6df93935ca6bc7d858b941fed300ec2104d698f6828262a7799db477a0962`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5e103268adbb270a22b1814507b65ff4e463061802b636036a208746b72072a2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef613a3f2e0247a22e15972035d60cd9208834d79b49ce19499d148fe2d4c748`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `67955d588c1720af190f340f719e936447e4511baf3df0fdc04f1f06ea411711`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f4e980bc89b4368defe54c1cec394fe99ef9acb064d0a0018c5aa191a8382b9f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ae56487a4b13c8a116b6a2721b7e45b3f7a63a82637310737886a4c1ebdc909`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c1dfc650c1fe23b9d042421b30de6c3dc56b6a5003b5c5bd9fa473bda77c65f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6695781886e679a15c9d86b7dcac575f04051026f212840d27fdafd066ecabcc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7e56d03475c38fdc4fac0e8bccc95e5274578d05f6ea0fb4832d4f670c2a99db`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `12103a2e0519649f496f7c3cf5e6572e2a5cc84d684bea57b7477a19d3ac2d18`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6a41694b4344f98e7ebdf3dc7ab7a749523b1cb9e420f8f356a93aad74b54baf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3d0440a13adfe09c180146b73d1d6978c2fb887094b2665dc7e3371f68014d9e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `890d9ad34fe1bac514fe950770ad1f9b4ba1154960517607460a938c6ff3aeed`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7e14d1edb4c244786bb2d64daf85f56aa133928df21afa9cf7dd496cd9d11ec1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5f9cb0b999f697b09c22465bae593614f7f45e004e62771a43cb4032a0ec8155`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f880cc542f85c3af70894c98f4a5103ea00ecd8856628457cc3e9cda7d216d68`.

The model selects the pinned bands at 2, 2.8571, 4.5454 and 8, each inclusive on its lower endpoint. It preserves the final subtraction of one eighth and division by the input. Positive infinity produces negative zero, while smaller-than-two and NaN inputs panic through the nil coefficient table. Lean proves the excluded-domain result and the exact infinity word. [Model](../../../lean/GoLean/Math/Qzero.lean) · [Rational kernel](../../../lean/GoLean/Math/BesselRational.lean)

Pinned compiler output confirms zero fused operations on amd64 v1, contraction of the denominator's leading addition on amd64 v3, and eleven fused polynomial steps on arm64. The direct darwin/arm64 overlay matched 66 boundary and 60,000 seeded inputs at seeds 17, 2718 and 65537, including 10,475 panics. The maximum error against the independent complementary Bessel amplitude reference was one binary64 step. [Comparison](../../../difftest/math/qzero/check.go)

The independent control used source decimal coefficients with explicit split arithmetic and math.FMA for the three profiles. It matched 44 native and 132 Lean observations, including the adjacent words at every band boundary and the squared-overflow threshold. Removing denominator or Horner contraction failed eight inputs each. Swapping coefficient bands failed four, and canonicalizing the negative zero at infinity failed one. Native amd64 execution and proof of the finite approximation error remain open. [Reviewer control](../../../difftest/math/qzero/review/main.go)
