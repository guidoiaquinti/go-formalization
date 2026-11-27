# Independent review of `math.sinPi`

**Verdict: FAITHFUL**

Evidence SHA-256: `d0d6701c52823d6d724c1576adc3f77b6ca795641aa15f983b10588693d67998`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cb13f0a7f7ff797276f8c8fe9bcf2a8925b723bfc7c4aa7c5197d666e370dba9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3ce458eb2f71dc7d2608e1e5880af488541a0c53062fc8982ac50a5c9c70724`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `209b86a49c658941577ddf65d181dd9d8cb78c2e869f2a92543dc4b22a18e644`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5824d9737fef58248ad441856ec1c9023f45a8bb68be8cd6730a5c2f4e036c04`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `51e13f8304ffe7ea14de49971b1ff8f926d127138d696868cc9e261b1d5fea19`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `075b57a651d5a423f8ea091120eb62cdba648bd654dd08356c732a88d9358e9a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fa62e7f57ccca7f21ee31c34cacd1df3a9c537ee8046bf5502048cb09d1ce565`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ec4da08072477d9dbd9a25509ea4db36c7a0c1121ef7e4fb49dd551a9aa71e4a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7071ca14d8186e7ee37b1c036a522bb2c9de42242d4b9d169b7b1ee6e4abbaba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8e2f6c3482b09b11eda170bab7c1592121398786ccb000859856ab14727172fc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4288ae8e5b287a0fdd4c7d7bd5babe490b0db9cfdf524270afa407a2c37e3c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2c7b8bb8aac6a1fe93f7fa7bda59b71c5c1ee223c262b840ae79c55d833245ad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bd596d4905947052d0c4a087ea0a1ac98e3c64343c5feaec5ea19fdcec4fab59`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `34a2a16c1a38a1e75f433aba9af158c4b9051c3caca29c7a6c64b84d3893f293`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `610ffe1fc5e4097e8f2ca9e566eec775ea17ddcaa3c906af646def141262d578`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2bc7ba6f41fd7814f16e2fa1b522b8515dc61bd119aa9da029033c998664b82b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0710961e834f70819fc9b10907c6b046da588de6a4af5e0184c40d43dd613765`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1c2a56939cf365153c9f29b0ba446df66655ef2babec01c290484d62bc181c1a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `423113e4800adc7b5b6205d8d1cd8b0104611cd45993da79e0fcbd579b702d60`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7f12313f5682f2a2fc862ad660702d1b172b68d5147d55d8caef659f230e2640`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `38a1baed4199dd7a2c8c3838a186c694cfd01fd7e8c7ba9aa6b86f39450189ad`.

Pinned Go 1.27.0 on darwin/arm64 computes the negative sine of pi times the supplied nonnegative magnitude. Its reduction uses modulo two for fractional values and binary64 integer parity near 2^52 and 2^53. The model preserves the direct small-input path, signed zeros, and quadrant reconstruction. Inputs outside the Lgamma magnitude domain are outside scope.

The author comparison matched 48 boundaries and 60,000 seeded cases (17, 2718 and 65537). The independent control in difftest/math/sinpi/review passed 17 native fixtures and 51 observations across the three modeled profiles. It used exact signed zeros, unit and subnormal values, and square-root values at odd quarter units, with two encoding steps allowed for the latter. It checked the representable predecessor of 2^52, integer parity near both cutoffs and the largest finite input. Omitting reduction failed 24 profile observations and discarding the output sign failed 36. Native amd64 remains untested locally.
