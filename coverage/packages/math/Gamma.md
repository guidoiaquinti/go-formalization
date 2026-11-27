# Independent review of `math.Gamma`

**Verdict: FAITHFUL**

Evidence SHA-256: `61a67dedf77adf4f62d4c4c38a9170b93524d906b57ea508a2b884cb4a76fbad`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `389f5054101b9ca42ae20138498f1a4ac9cb03d8b8f0ed1ae3025d23bacfb907`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b4664898d63ce729d0f3bdc7efc4d17cbe55ffa0c1467e2f871115d7092c1183`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3f805b37f49a7d5311082236f80e9d5d2d62801786e264537dd38b8189c976d9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eacf9d718e4df1b21efb8fc913632be4e0fc797faa1355d19c79ef22f3a4dad9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `006a24d3bc04b18992548eabf29646c460fd40c03796cc49971ad6a31615f98f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dee53ac495ee9535d362eec6d627710f36881f8f822d9c75f4aecb3a5f110ace`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d93e2d29635ff052ce28dc539d1cf33b1014f8e8b55560a0116efe4fb97b2057`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c7b3b160c9786d2d2730a84c0cee413a89522a5b211f1f74663eb5ac1900dc38`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7a3c3d3b05f926f973b35422c6e8c5574fd83d450a0de23eb65e4ef78d810324`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eecdf3d778559e15f97911c3a276ce8e9d5702b9b28e16932b9e015492d3108a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b3b75f6e92663af3732984ff4b2759c9922256289c635b40dc3018f1e7e24b3e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `25cb184a543445fc84871d1cc16e38a1426ee74cecd3ded93c25587cc90c0fc4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `29b4a2589af52d7c6e5c29cd8fe67723626568fe5b38954275943aa4c1c9fd3f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dcc9472c4e9e35637b745090a63aeef78376a38e34cd0d05aee93d2f1b6e41d6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `145839b8445a53a1bd4e6ccd5fb8b50f9e311a3b876539ab9c76a9faf53db232`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1d5bee4083d433624e3af93f7cc7a63f43d7ad5ac28553d498809593ec9f7050`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4449778662bb39ca2e7b7dfc720d9e48f5287cd7c8a083b9de16d8654dd9c46e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fda4d6068b8a43f719d5a0c4854e46e9f6b97798c631b6548ce5cfcbefd181d4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad7e115f18f4b5f67899de125f9d6aa7e7a26a2d20417acc5e3e7f32d62a42b8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76ac288e2a3cf91f528762843c790aa17bdf38d28729cb71c80bb3e0f957ac1b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9ecb8fbdb00b95ee701298bb30faa11e464fd1a8b7b236a7400b5a3ebb637ced`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4abaa4cd8057d062ad2266736fdd76d09a65c82c31b156b38d90cf3c65cc87cb`.

I checked the Go 1.27.0 recurrence through magnitude 33 and the Stirling and reflection branches above it. The model preserves the sign inferred from the integral part in reflection and the division fallback when the product of the sine factor and the two Stirling factors overflows. Signed-zero poles, infinities and negative-integer poles match the registered classifications.

The pinned darwin/arm64 comparison passed 138 boundaries and 60,000 inputs from seeds 17, 2718 and 65537. The numerical reference differed by at most 1.09e-14 relative for normal results and 25 encoding steps for subnormals. The independent control passed 59 native values and 295 Lean profile rows, using factorials, high-precision half-integer recurrence, and local numerical bounds near zero and the cutoff. Discarding signs failed 13 cases, omitting the central recurrence failed three, and producing zero when the reflected denominator overflows failed four.

Pinned compiler output showed fused Horner steps on darwin/arm64 and separate steps on linux/amd64 v1 and v3. The tiny Euler correction contracts on arm64 and amd64 v3. Native amd64 execution remains unverified locally.

```sh
make verify TARGET=math
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/math/gamma/review
```

[Model](../../../lean/GoLean/Math/Gamma.lean) · [Comparison](../../../difftest/math/gamma/check.go) · [Independent control](../../../difftest/math/gamma/review/main.go)
