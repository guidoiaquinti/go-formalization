# Independent review of `bytes.Buffer.Truncate`

**Verdict: FAITHFUL**

Evidence SHA-256: `6cf2067a9af291b0b353358ac69e8e2d12307e42f93fae786b0f41eece707790`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9932e74e3658f8a2d5083a207a1def92b8c38b295a5fb772645c5eb53987cbe9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0a587ece18dc0e74bbe6e83dbfa138298490b7389435c5d3eeb4cbc22ee04185`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `85c0a8efa25125626aa323757ccc5b7cdeb0cd6aec3557ff65045796efebc671`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3fd042aaab7af5113eafc21a45ba385695f5cad77b049a2b693cb6bec863c413`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aa35b198fb2f08dade75cc8ec27ab8832b9d3cdab970cffeed12a4a8e68c81ac`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6833d7eef64bba2383e2b34aa3fe7b30868d6648289dbeaec857f159b681e72`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `42fb6dab825db8defef907ab355eb791371cc796ae866a9f7c4862bba1324e71`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d789d2dbdec2b6329d0dd2878344f4dd9bb41bcfbea12bf4122ade6cb9f220bd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d1e740c5be83a5a9e089024afb900826d12aa4e1e642d80c5d067f855073d334`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b7ad15547dcf0055df9076c02328ce6d4ee27db8c3b27c2fc4f69e9acb47a173`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ab517ec1f8d4cebfdcab1a6b425c0a1ebcb8b44cff242cab08f000614cd233a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `00522c9d439495ca2df99b5396dbd28faaee8afd06bf5e5504f53d7e49a6d6aa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3e7ddd1f2559d8594673fbdd722ed56987e16404fe012a57cce6f90b5db10b67`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3b2fd611c4e89302faeee0a149432fdc0347d97aa3d754905899db2eb636b06`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `619a188aa0e637c70bdb29612c8582317bd65c5b2b99e7d85159ec669f4381a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `292a02322c9209500c1bdb5f6e76f5584ccb35d9fd02a2eedc0aaba5f9f39670`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `61830542ef690f9ef1e2570f8fb8b4abd21de33e4f480a1a98f370c76e81d441`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `044811e1b25ff0f47d0b34999528f238a2ed6b38775994afcff4d29dc8e2ccfb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7181e391c80f8a9edfae4535497968145db19c020d79935977414899c228854c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dc13a6cd3484526fffb43aa65cb18b1c363d55532a4a95feff329cb6c5e95ba6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc0fbfbc97d2ce439d422b6fc6278d1c66b0e6a0c1e7884f82ae11f84a9af810`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `72407a67703a06601827f210bd7444a426726393264d1db5191a05d6c14f7de4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2a2c73529a833d63d3f00c89c5f34d1444db949dd1f79abeea5e1f98db559b26`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `07a4c250518f927ecc4c3b07e29730b8a140c666aa4671f455ce9d7c7e9fb187`.

The pinned method delegates a zero count to Reset. A positive in-range count reslices at the read offset plus the requested length. Nonzero counts invalidate the read flag first, including negative and oversize requests that then panic without changing the backing storage or offset. The model distinguishes bounds panics from nil-pointer panics and proves both state transitions.

The darwin/arm64 Go and Lean comparison passed 11,056 boundaries and 60,000 sampled inputs from seeds 17, 2718 and 65537. The independent control passed 14 fixed observations, including signed integer extremes, ordinary and rune reads, nil storage, and nil receivers. Delaying byte-unread invalidation until after validation failed four bounds cases, and delaying rune-unread invalidation failed two. Reslicing at the existing offset for a zero count failed two cases. Both zero cases reused the original backing at offset zero on a fitting write.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/bytes/buffer-truncate/review
```

[Model](../../../lean/GoLean/Bytes/BufferTruncate.lean) · [Comparison](../../../difftest/bytes/buffer-truncate/check.go) · [Independent control](../../../difftest/bytes/buffer-truncate/review/main.go)
