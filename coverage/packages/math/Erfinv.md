# Independent review of `math.Erfinv`

**Verdict: FAITHFUL**

Evidence SHA-256: `97b46ebab05cc99209709f9de3870d289b39c56df465c221706f85d02af1ddf3`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `59394e106628ed5529ab7b68f148dfd9c3308cf3501bd161826b41a3df4e796b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `26bbd88f77df0b56e68c7a068ce3024e53d880c809eb4d4c6b4b1ef095f97594`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a226f3e4a182d2a905f47f5f7284e03cd29d5c443470dfdbdb5b4cfabfe90d52`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `62430e040c54ffa1767bc3bd3786a288a09271f139bd468568af0e4d12d720c1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `55a5097c636ffb87bec05e5343b74a15651dffd3056b30881ac66bdf594f308b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4346ea85f69cbb811481837454b71efb813b0498fecd90be54b111d953f703a9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d4010364fd821d19fd10f678dcdefc8f373eca41d9d5d45083aec9730dbfa5f2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a72ba7dbb3e8f0e9a911873c076fefd2b9e6118c10609190478fdb4645dbd906`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `66c7f799dbca88e46f18aca7aed9b3145dc447dab54e7e1ef0d5b1811457af63`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74f596336e9d47143559d14166c15435ef3e149236dc8e7bdf7812bb689f1695`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0862ae249fc5b031da39fb614104364cbd9445048ace2c8403591590ef5df4f1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `925774c9474476688f0546c464f8a8a38bca4538d6076d1848712b8ad39a0078`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3e3b4cf141021b18fdc396412cbdf81eda2b5bce694c3132e77f6c51905b72d6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `03cb81741cc24304a82cbdc6b45073f4191cbaa21676b40b828c593331e718e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `00314b50a8767a1d0022fa1b1e22428a8e17da5c0ca416cc16ce04f9e204c82c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `58f7e7138f8890309d5d94be24e3ed070b46191fdaadccc0fcd3df602d7c337f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74143c8e058ca912c3681bad028a5159af3e2a9fca9778b1420fca5a95043ec7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `496a4ee42c0a3dda0271c5134d402fe1e5ba6f6516ff50805ce9018f569c6cda`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `636a4a29904826cee88598f7c195c301304dda0924e6bddaad8953e98c487b4c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `577c2ca71e0b0956f8d903a36dbb9b63f282d6da6b9e075ab991219a094515d8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `955d64fbb35d9b7562e588cda06cdbfbf5244fe38a56df608f8a32b856f77500`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `79c6c77649fa815a51e402c993ab219f8ea1e6077740e0a49ce4bdb3304e5b44`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d135adeb123e794c7d45dde22cf397279a38c8ec769b427b60cda7471b20cd5e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `037fd5015298e1fe5d296605a3cf3deba60c563340991510b449cb36d06123fd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ee8d79a6b099a31228a61ce6c2a9b905203faa88c93b031a6605325cd5e14810`.

I matched the Go 1.27.0 central interval and its rounded radicand to the Lean model, including the Horner evaluation order. The tail uses the rounded subtraction, Log, Sqrt and shift for each side of the secondary threshold. The model preserves the negative zero and signed unit endpoints, and recognizes outside inputs and NaNs by class.

The pinned darwin/arm64 comparison checked 78 boundaries and 60,000 words from seeds 17, 2718 and 65537. Exact comparison agreed, and the registered numerical check bracketed the interior inverse within ten encoding steps.

The independent control used a 384-bit positive-term expression for the forward Gaussian integral. It bisected binary64 encodings to enclose the mathematical inverse, then tested 54 native points and 162 Lean profile rows near the central, tail and unit boundaries. Dropping the sign failed 24 fixtures, using the linear central approximation failed ten, and clamping before the unit endpoint failed eight. Native amd64 and other architecture profiles were not executed locally.

```sh
make verify TARGET=math
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/math/erfinv/review
```

[Model](../../../lean/GoLean/Math/Erfinv.lean) · [Comparison](../../../difftest/math/erfinv/check.go) · [Independent control](../../../difftest/math/erfinv/review/main.go)
