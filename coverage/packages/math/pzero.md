# Independent review of `math.pzero`

**Verdict: FAITHFUL**

Evidence SHA-256: `09882881cfc2a9fe33d6821671c06c3a64d8ad1b68b655de85de37f24c07853b`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `967e173d3476daf3b98625e9e6747f86185a093fd62ef0e3d1761baf50d7e952`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a038016b73ac76cb80be20dd217f610e08daa9382ef99ed0b97bc8183c18e217`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cc5aa9c5595fba9478a8d7197069c74c31c5a9e9cb039fd05ca14541123f85d1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad9a4f30bbae7910d95e1eab8e63b3ca7f2f0be47462e5859f8d1a285303d587`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2de6cc2b0a6888eb427d7bc7c8fe5830ee5309fbca71618d8146c6b98f65bf6b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cca4db89e9a8d51fb45d77071151f8e69322b77deba7accdd1eec21fde022f17`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `16c8795bd37cc6bfb455d4802fd3ead6c339c4b2247a76fbc9a6d5582d10cb28`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `127dc266f4598ec689009611fee4ea0b2922c4cd47711340036f8e2a9f524688`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d5c2e5de499232e3560442012bff63d4f5ef1a727ef34affb26cd7217fd17fdf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `99594511a2d891187f2e7a9fb946e2dca422d05139392add1eb3cfc1ce69a0a2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cb617b6b4b3c47d0b432691dbfe8a1e83504a84fb22bb6bd2e532a8659a308b1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2517be65db9977166c8846fa20e59f637ee25c39c840e5868dbce6638a65986a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2be62b7795cc24073e1c50ead1c8cdd72d4671b2511602aff1d84a76d7f40472`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aaf856f606f49dd4277a8e130689f4b7109da53c7ffc38f7887e97227963e940`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fc152b8a7fa757890059650d05b2dad6f29c566a416fe55cc2a9009028efb48a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b3a6de545b68b5e22eb29dc95366c63a06893b6909444503d20a29da032d0290`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3a62cf8cbde175187a8402839036c7dda0cc130641d9eb6f4bc07a5bb2d7ef51`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `26a703aedf5683c6c4c37201e41a590d16c8179c8b5d974a22d865f0ddf1b9b0`.

The implementation selects the four pinned coefficient bands at 2, 2.8571, 4.5454 and 8, using inclusive lower endpoints. The hexadecimal Lean thresholds match the binary64 encodings of the Go constants. For inputs below 2 and all NaN encodings, Go reaches the nil coefficient dereference. At positive infinity and above the squaring overflow boundary, the reciprocal square becomes zero and the result is exactly one. The Lean theorems cover the excluded-input panic and positive-infinity word. [Model](../../../lean/GoLean/Math/Pzero.lean) · [Rational kernel](../../../lean/GoLean/Math/BesselRational.lean)

Pinned compiler output confirms zero fused stages on amd64 v1, one fused addition of the denominator's leading one on amd64 v3, and ten fused stages across both polynomials on arm64. The direct darwin/arm64 overlay matched 66 boundaries and 60,000 samples from seeds 17, 2718 and 65537, including 10,475 panics. The largest difference from the independent Bessel amplitude reference was one binary64 step. [Comparison](../../../difftest/math/pzero/check.go)

The independent control evaluated source decimal coefficients with explicit rounding or math.FMA as required by each profile. It matched 44 native and 132 Lean observations, covering adjacent words at all coefficient cutoffs, the lower-domain panic, positive infinity, signed NaNs, and squaring overflow. Omitting denominator contraction, omitting Horner contraction, and swapping coefficient bands failed eight, eight and four inputs respectively. Native amd64 execution and proof of the finite-reference approximation error remain open. [Reviewer control](../../../difftest/math/pzero/review/main.go)
