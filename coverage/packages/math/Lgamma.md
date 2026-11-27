# Independent review of `math.Lgamma`

**Verdict: FAITHFUL**

Evidence SHA-256: `5f6a484fdf194676dc1bc815b6e85006a5a2fec3cc7a0a0b129a66ecb76fb2b5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d3c49202e9c06724f57d77b0029ed19ef645841989ca3fd05d93b25664b80fb3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ceaed7fd304355161f7848de97fcf285af926309e3a234b1152a66bb96c834f4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4c1a368acee4677f704138f9d9ae47dd6bb2cb97993dd18153af596954618892`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4b776a17e5935e7aa0322b761b6e864153f16fb68a4ff8aae3675631614f15df`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7f174359d56b08009660e03798bc707746e9cb821f4b995099971b250d5b98e9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4b82dce6c38e5f3bcba77c83a06a6402cbd2023eb867e453c57f8c8f81858c63`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `041224c538bd193eca4bde91a715f110c39c2f3cee481cdfd11d1500ae094a85`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `525f8e81e31b19c66b4913c541f7aac4c8657f9c1d86aa148fc71e4717c918a2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `484c18eefa443ea9dfe45867850cae2c484fdb7c65013fdd443e566b1a60ccfe`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `990b4cb86766a92f52bbe0cd9f8108abb4d1613505895cdef15d3a69230e306e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `09742ac6fec911d0a48cb8278573fdff087ce8639592f9b3c7ed0808130c3f57`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e7348d44e0dbdc1ce2501409276fb36da70256e3cb9a5efb5f8daf7f3096f93e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8999d58402eda658d412683ca4e9720b5ed30cf1d7647cac07e3622eb7b2350c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `99b1732861887233c0ac232fed680fc5fae471c7fc097ccfcb7f3f8736129486`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `070c4ca37c112aef2c9dad410bce2129d381d9051922bdb39658aabbd4abf481`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b6838a39d9e6b7827d8cf3df5943ada5597865414ef5d0523a57f55ace2bac9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cdb37b654b92c7a40ea3905ba8e22ff35f4bc9d97155888d976d4619f3cf5bf7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e37736d13e2b642218e85a6913b9cd91af76b7437fdedb128cb7bc9cecd809c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8efc88abc89ea43693b15c58f827f875207440609acba2992192280875e170a3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6436624212dc8a8f6483eb0f527745f3face9e300272b6aec50be00bc89af9ef`.

I compared the model with the pinned Go 1.27.0 implementation, including its poles, special values, three central approximations, recurrence product, negative reflection and overflow. The native darwin/arm64 differential check matched 144 boundary inputs and 60,000 samples at seeds 17, 2718 and 65537. The normalized error against the 256-bit reference peaked at approximately 5.15e-15, below the 2^-42 bound. The model proves NaN observations, signed-zero poles and the negative finite poles from magnitude 2^52 onward. [Model](../../../lean/GoLean/Math/Lgamma.lean) · [Comparison](../../../difftest/math/lgamma/check.go)

Pinned compiler output for darwin/arm64 and linux/amd64 v1 and v3 confirmed the modeled contraction sites. The tiny amd64 branch uses the public assembly Log. Its subnormal arguments acquire a leading mantissa bit at exponent -1023. The corrected comparison checks the effective operand at high precision on amd64 and tests both signs at both subnormal endpoints through twelve synthetic Lean observations. Native amd64 execution and NaN payload equality remain outside the verified scope.

The independent control passed 58 native Go cases and their 174 Lean observations across all profiles. It used factorial and half-integer identities, recurrence near negative poles, normalized logarithms for tiny values, and the leading Stirling expression for large values. The four incorrect variants failed 11, 13, 4 and 2 cases, covering unconditional positive sign, logarithm of rounded Gamma, premature infinity and the amd64 subnormal stage. [Reviewer control](../../../difftest/math/lgamma/review/main.go)
