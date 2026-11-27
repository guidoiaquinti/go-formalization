# Independent review of `math.Erfcinv`

**Verdict: FAITHFUL**

Evidence SHA-256: `69be0fc757c108e0474ac3d7165855e2d8eb1c57b661f43eb5c0292ce20b4f5c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `483556f7ed72b4da5b891444ff785a80297a21fb3268dc380a25f6acf6eb0460`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93adfe82dda40f2451ca68f7eaefb599b020c4b1f96aadf0501734ccd702c706`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `96a13400818e1de777db79b9799950b896f4f03c056a251b5a8bea966ffad4ec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `254f85711c1af3cf5b9544ca397d6cfe0917dea6f9df26208468e920215c08ba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `17d5b933828915e5fa22ecb665e75dcbe5b2d16e971976d6254a774c9522f326`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7f606b2e7ef94850d31fc3bddbf55211a5e9d83b05c986ea2fc1cf982e392c7a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d3554608c184407a4dab7dfcb8e9b14fd69932205e89e81111128af0c9989cc4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef562205e54d99f8b25867478ff4ad9598beda95a0472feb7e1142340f066972`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3a4ccd69dbc9e3a4354a7c4d44b8e1229310479ef89283823e8717d4868ae8c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `717416311aab35dcddf1e629157afb992aac986181ea75b175eb825a1467d458`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `14bc869c8c01ae3cc7bc3af6a28430d39ace565cf2c999f3be2d0216827b5de3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dc00c9ef17530063f8661a154f759091d171614a5d7211ed4b8d8b8d25f60349`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2f51e73961fd4ffb53c1d99ac76e1610efed37e38057356b6c61bedb96097df2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b3403a5bbf472618676330b43774f218bd37719fbd3ec2db0f7f49de3b91d45`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `45911d458464e5116c13828716c08742f41c1a69f337d9fb4675979168bac00d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cd0887a5a002715f8dd9ed0f0e82ffc2661f590a6e1f6f7134db2773d9834ee1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bf7aee31db6854205051b6666f8b3799444f0e040caddc735ec46b1dff6d1d6c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d01687dafa176223685290b57d359e411004ae1e7cce2b065b6cd090e88d0503`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aa63e16f781a6051c6778ab2b5a0e7423247e897fccef725ecf2051d413f01d2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6331cf07790671888ee7fa33ee8cde79efae1fa0e0ae5f23acd0a4d3a1aaeac`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ec415ea212618db97708e9d00907bbad9f7b80243db3c4ac75145ff6b9e8bbc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7708dfee0b33efd3d101d4b7495336e3c165b4e32c175a075e8770a1a212103c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `873bf9ecec7cde5bd747c16775ef2d722438b48ad8f89219b588aa5a453accd8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2f17234cc7063c74d37ec1c48a0fd2cc16a53a16f3f00ca014153a9df6ccaa09`.

The pinned wrapper computes `1-x` in binary64 before calling `Erfinv`. The model preserves that rounding step and the inherited architecture-specific Erfinv arithmetic. Lean proves the nonfinite class, endpoint words, and concrete negative cancellation cases. On darwin/arm64, both the smallest negative subnormal and `-0x1p-53` round the inner argument to one and return positive infinity despite the documented negative-domain result. The next more negative word produces a staged value greater than one and returns NaN. The analogous finite-precision loss on the positive side extends through `0x1p-54`. The numerical accuracy contract for that positive range remains an open finding. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/math/erfinv.go#L119-L129) · [Lean model](../../../lean/GoLean/Math/Erfcinv.lean) · [Findings](../../../FINDINGS.md)

The native darwin/arm64 comparison passed 51 boundaries and 60,000 sampled words with seeds 17, 2718, and 65537. The driver independently rounds a 2,048-bit subtraction and compares the resulting stage with native subtraction before using the separately checked public Erfinv result. NaN payloads are outside scope. AMD64 profiles were modeled and native AMD64 execution remains unverified locally. [Differential comparison](../../../difftest/math/erfcinv/check.go)

All thirteen reviewer fixtures passed on pinned Go and Lean. An exact-rational subtraction rounded back to binary64 independently confirmed the fixed stage words. Rejecting every negative input before rounding failed two fixtures. Narrowing the staged subtraction to binary32 failed three. [Reviewer control](../../../difftest/math/erfcinv/review/main.go)

```sh
make verify TARGET=math
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/math/erfcinv/review/main.go
```

## Go-only contract check

The retrospective probe checks 11 deterministic inputs against the documented special cases, including signed zero, nonfinite inputs, and values near the negative cancellation point and upper domain endpoint. On Go 1.27.0 darwin/arm64, two negative inputs returned positive infinity where the docs specify NaN. The next more negative word returned NaN. The Lean proofs establish these three witness results across the modeled arithmetic profiles, along with endpoint results and the NaN result class for every nonfinite input. [Go probe](../../../difftest/math/contracts/check.go) · [Lean proofs](../../../lean/GoLean/Math/Erfcinv.lean)
