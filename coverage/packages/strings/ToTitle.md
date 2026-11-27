# Independent review of `strings.ToTitle`

**Verdict: FAITHFUL**

Evidence SHA-256: `b8c9a97d6d151048aee904a7a0d3edd2d355c9950764525ef093efe5d7a9c2f9`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `569648a4eda39115399d91a7cd8de136f831c54bb47bb303735b7bb700d449c3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a67ae5661b047ae4f996d5029adf983d721958f79bb9330f0e400c7f64d8667a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `294292fc51ccf1f5ce8a790e947319fa09775026d587bb35f938b0ef43d489c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `baf113a3e13fc40fc9f3b3564a37bc00eb1fd426c6fa8d1cec2b73fbdb3fb4c3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `489714bf88b453d6f5901ec85aeb34a0b27104c92b121fe03fae3cf1f95c61bb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `10a391e76d41d31584c345eb05f3aaa5682a7c050fc8d146ccd723864005bf81`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d659987370b7a78479757db96120966463c64ffe860f182b580f77f549b91a2e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2aa2ce1cfd321709aa32fd5dc10c2038508eb336ae2f79b880aae72f20124f48`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c13ee54fdafa7f9150be7b74b591f30df63bb9b86ebe7f21c1323adfe54dcc79`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8fb5cce41d5b492cd497dde1b5e4b1e129dd11c445b255e36f16a7f44cff22af`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b97d0331af8dddde0d4fd79edc42cad8994c2e18276234be5e58194384b401af`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8105726640025ddeb26063df759a289a815431e29896b23e95fd4d0af94e315c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `85a50f183d8aa44e31bb7677f88a69be0977eeb5c103580a06aed0e634430727`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dcd631a7ceb8ce7d0f4d2ff741f89b1978491f93faa4424904e6f84a6c71292d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5bb39992af432aaf2e1a1d04cc1ab5406536868761105dee3323f976c2dd5dc2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `944d51fd6b86d6630bd90ae85d6222e4e1539f1ae6deb0fa56eddfc45a05fc5b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc2f3547e6dbe9c74818b58a9762c44893541e19faef0f8368ca71f17ae30e43`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `073e8237137f35e291a1089913becae211c3d764f26b8c06782de5de6eb4bc0b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `82ca8e35625736d9b9b5a1e250282c0c455511b899036c4dedd40e3925774826`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e751883738787a1bfc6bf0f79e6e29027f1698831d88ddd2d17af64f250b4e46`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `81903e8e5bd96844fa3a395612344da375626b699c8f2f8d144ad29496d4bd97`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cf20574868b55d7315289bd079d21e925cb8905d9939c8f6a27d32f454f0536f`.

The pinned wrapper applies Unicode simple titlecase independently to each decoded rune through `strings.Map`. This includes ASCII letters at every position, and each malformed byte becomes its own encoded RuneError. Lean proves the Map delegation, empty case, and an ASCII example. Allocation and backing identity are outside scope. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L747-L749) · [Lean model](../../../lean/GoLean/Strings/ToTitle.lean)

Regeneration with pinned Go produced a byte-identical Unicode 17.0.0 titlecase table. The comparison passed all 1,114,112 code-point lookups, 3,446 boundary strings, and 60,000 samples with seeds 17, 2718, and 65537. [Table generator](../../../difftest/strings/to-title/gen/main.go) · [Differential comparison](../../../difftest/strings/to-title/check.go)

Pinned Go and Lean passed ten reviewer fixtures. An ASCII-only implementation failed seven. Using uppercase instead failed on the distinction between titlecase and uppercase digraphs. Collapsing adjacent malformed bytes failed two. [Reviewer control](../../../difftest/strings/to-title/review/main.go)

```sh
make verify TARGET=strings
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/strings/to-title/review/main.go
```
