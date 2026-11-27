# Independent review of `strings.ToUpper`

**Verdict: FAITHFUL**

Evidence SHA-256: `20f055de0fff9c216c75591021fca2bb9f372b6d9fade265c42986c7881eca18`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7abc7dcd17911f09c2aef7f2021efbdb994953462612102c568d57e40b3a5d8e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a4741e7bd6028f4c30dc76528494b258c3ef05ec23bc5959a9940fdad95cfa0b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `48f3a6349b20e77dcec1fe828b4c7ea3c5d732e2889a9ed5c3105fddd6ac484b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a6b1985499dd5db8d69d1d75b1b7091f715d4dcccac7448a7975a30017ccc68b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a69074f33e7d022eb78d9572c58abc40f44203d3d8d4ff402dd77a4d5f4ccf26`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a914e4b98e466497d4f95f13ddd06cd4b8a6bca0f378d8e00331f887c958b5b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `54aef79ea8f456b9e50aadaed28c8568d57f0a828e0a344d74f6c3d686a9334e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4f519e25781fe1d634c9f7f50e8f105d938b56a2c7b8b6d4825c62e3f7188f73`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f07bba97be71f601d6cb0de76be5ea47512be3dc572ae694b9f6f09b04d628a6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4ba9b7be06bd5578aff605c8c2529fbd998d836572e2938740e69e25c1fa13f3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad4fb448873811e804c455721d1eb18d0bffde3e67cc687d7a0cfde58707b823`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3861b568c62f7aeb71518c97b52b47eaca634df43d786bb5f75701ab5e858ddc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df14915f9624f14492832a9203d1d9c53522ce9a2fdc3c242e44b0f12066cb27`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `219f3011a4a9f07b30d7fa2bc05c0f30c323dfa6343480c93095226a3cfb3d7f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7fecb90f742d33faddc63f1bc670f5cea0fdeda0ef5dc85ef1fdead6c6451613`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fe84e9d9e5960adf9e7e7c184de0e3f60467b9a5645dce31b9d501e674dbf3f1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90c0afca72609b4fb8fd8735570652629b6d86baa0e386c528c70bce8bcd1748`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8c471a475ce8fc540fb11278e757df4ef999230b3e4b42c5083d7b939f61b09f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1745ef46d6429a78cd0e4576c69b944bb6dd8084a8345ce62fed8adf6211562f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `233913ac2abc1b97da55969536338f4345e5ddda178bd7ed4ac6125174e8c2ab`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a956a7e05df1541fbe62e6c0abeb9ea0c3c1ee9f6892d1167c7aea2ac2460159`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `33b315b323285602798f052f117f94c8180714b304b5c2d05acde8e36f5ee1bb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `41974ad36ca47286fe235f15e5444acbc54edb3d6116967908d36fafc5d4b088`.

The pinned implementation uppercases lowercase ASCII bytes on its ASCII path and uses Unicode simple uppercase mapping when any non-ASCII byte occurs. The Unicode branch decodes each invalid byte as RuneError and emits its replacement encoding. Lean proves exact ASCII conversion and the empty case. String allocation and backing identity are outside scope. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L668-L705) · [Lean model](../../../lean/GoLean/Strings/ToUpper.lean)

Regeneration with pinned Go produced a byte-identical Unicode 17.0.0 uppercase table. The native and Lean comparison passed all 1,114,112 code-point lookups, 3,538 boundary strings, and 60,000 samples using seeds 17, 2718, and 65537. The string oracle applies Unicode uppercase to runes obtained by Go range. [Table generator](../../../difftest/strings/to-upper/gen/main.go) · [Differential comparison](../../../difftest/strings/to-upper/check.go)

Pinned Go and Lean passed ten reviewer fixtures, including final sigma, dotless i, supplementary-plane Deseret, sharp s, and distinct RuneError outputs for adjacent malformed bytes. An ASCII-only implementation failed six fixtures, expanding sharp s failed one, and collapsing contiguous malformed bytes failed two. [Reviewer control](../../../difftest/strings/to-upper/review/main.go)

```sh
make verify TARGET=strings
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/strings/to-upper/review/main.go
```
