# Independent review of `path/filepath.cleanGlobPathWindows`

**Verdict: FAITHFUL**

Evidence SHA-256: `427ecfda0900944d75367adc4a3ea20f3d133b5b87aa454037862130065661f6`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b7ea38a4479c06b54a833edd97cafd8103496f37d1ced06666e9aae65c237034`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2a9502e39b2163806b1d51c80d3642364b81ed0efc890888f41e26281d9bed5b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4d8ebe1c7b8427a89998f9f2783c8e00a602b9c5053ace77f5649cebc671e23f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `85d385dd146231c5b9abaeec3ad152f5d74bffa28fef6ca3b25fdaebfec826e7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98b52a247bd80e3995dcac8718ab77a6e091a126132398b9d80b62345c1264c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c963af8d2b3f23a4ff1bf4a40e3a98f0c8263f1030ca6f8764e9a3f42c80526b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7f95e9f53c8a9f38ab070bf037d74d89ede0adf0e1d33c5e16f33da6d71f6fd9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8c34a2480c4e65f685c01237ffb49c04268e099969c15d7f089075801715ac0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `496de680475dd888ef2df5ecef46f830824228c4003170c3de1264cc40d27ad7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1e2a7e7e3d791c085f89dbf7982e82c9b5b8c9a42941caa35da3fb83e9503400`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `459ca2df06b3c674db8496e22177b7b0cd7ddb9ee69a6f5a448474b829192df7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dfe76f85208501b07759bdb208633461d2830550a3921e7c6d3456b1a418f283`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cfb617159a7a2e2c04f06aba9ddb8a746ad06e96341baacff6c1c0a65742b1f1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1e91e9bf9199ce10cf2759359819165c4a1f333f1fb919924c8504648ccf0943`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c02818381073017725c809caeb1796d01f66adf69bf720ade2ca498f381da08b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `241e93c865490db8a453d23083346b4007a70a25eef14670745c6dcd5fc43c39`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `801bfd75b7a569a8ee733cbfd3b34f48023eb1c1164ebb39a89561530ecbb81b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bb36e185d1cb99d813616fe93bfeef9e194649aff0650ba21ff9b8d001ecceab`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ced67f6ed1ab1c9ae3107a5009526f2ed6a93d8c445c2b562d6d18bcb9664cc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f239f00820fba56dda3530af4a5d3dbd143b4b617cb0569fa4a6ca92e3b1ef94`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01d73b8c1aaa756ee4e71e2f19114257a4d2036f9db5f0d865142c4e428bf132`.

The pinned helper uses the platform volume length and separator rules. Empty input produces a dot, a root retains its bytes and full prefix length, and a two-byte drive volume gains a trailing dot. The ordinary branch removes one raw byte and clamps the prefix length at the new path length. Lean proves examples for empty input, drive and root special cases, and UNC volume clamping. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/path/filepath/match.go#L297-L314) · [Lean model](../../../../lean/GoLean/Path/Filepath/CleanGlobPathWindows.lean)

The pinned private Go overlay on darwin/arm64, compiled Lean, and profile oracle agreed on 528 boundaries and 60,000 samples with seeds 17, 2718, and 65537. The Windows volume and separator rules were simulated locally. Native Windows execution is unverified. [Differential comparison](../../../../difftest/path/filepath/clean-glob-path-windows/check.go)

All eight reviewer fixtures passed through the native helper and both Lean profiles. Removing one byte for all nonempty paths failed four Windows cases. Skipping the volume-length clamp failed the UNC volume case. A malformed-byte witness checks that only the last raw byte is discarded. [Reviewer control](../../../../difftest/path/filepath/clean-glob-path-windows/review/main.go)

```sh
make verify TARGET=path/filepath
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/path/filepath/clean-glob-path-windows/review/main.go
```
