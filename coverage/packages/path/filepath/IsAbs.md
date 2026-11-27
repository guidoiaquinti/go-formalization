# Independent review of `path/filepath.IsAbs`

**Verdict: FAITHFUL**

Evidence SHA-256: `61de67f827c8248b9c1e67aaa8cae1f34b1869d95c05679227a4bb5bfecab1f2`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6517c9ecfc0ccc17f39c942353ec39d7d30b1e50dd23fd3f246af3c0623d953`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `35d39a6e52a15328ce1eff3ae659b9e1fcf0a6c30a2a9f2301f04bf1f6312ab3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9e83e1d3b3cfea313173437f7b059691c07224398cc0f918befec5459ceacf6f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01b04ca02bbd2daf9dab36caacdd85bc07e1937c777e792aafef4dbfa962125b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c7c4f0ed4add3a00f9fe88927382a2b1f607938db177fc56152cc667ffa52db2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cc25131d6717a9facf8d6ea015383722b1b396b4e343bda60854e2501fcbf0ba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a08aa7a6c88b7eab54eb5c5f456eec28e8ccb370595c4cc72f70c1e89afe3821`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4a26ad795ae4805984938a86ece0a301d2ec93c65fe1483a5db9e059b8c5c41a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `89c27f4c963e3801b85c576e34a9f1c55496ba6fbcb42081a6c5cf040851a7ca`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `58f2aac74708569821a96d8d0a953278067de0e1ab58e040cefa61e8e79669c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b1648ce8c352c79cd7f7f8cdd954cce4e78ae6ca362a1ce29d9f24b70b3b1682`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b702d0c3136a885a410585d9834d3f64665195af20912153d6ebc0ffa9833635`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8fa303f2e3afd9f2ff7491d7453eb2665fcbed81887360e3e167a20d40d8ac3c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3eb1f1a0c5fb1a659a858b61150bf13e011b89b11db21cb994f4fcf3c6e32253`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `49e17e79fa264759cd16f567ab98b1c82c6ca8efe512a25ad6e2bb357330a15b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `25b4cddf00e0e101b0f6c61f8e6018099bfe603f61984ed154b06cc76c2634a3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `689c1ec4a6b325f4a751b2d3724e389eedd2a8c741a26c372e7bc3d97a56f4ad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1b45ad5127fec16f3c70900d44479baa5be5a55efe51e9546d2dfd9c347db2f4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `41c0d26e4b73401f4bc6b8707a14db08f06a7f3c4692ec28004a917b0043871a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3c018edf1e30e05ee455da3789e22c321a9ebec064fc1d3cfb43905c1c42972b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3d2674c949443c0b9789adb397a5ce447639502776f98c7bb3509237b920c63`.

The pinned Unix profile requires a leading slash. Plan 9 also accepts a leading number sign. Windows first requires a recognized volume, then accepts a double-separator volume or a separator immediately after the volume. The model uses the pinned volume-name behavior and preserves drive-relative distinctions. Lean proves the Unix rule, Plan 9 device prefix, Windows rejection of a rooted path without volume, and examples of absolute and relative drive paths. [Pinned Unix source](https://github.com/golang/go/blob/go1.27.0/src/internal/filepathlite/path_unix.go#L35-L37) · [Pinned Plan 9 source](https://github.com/golang/go/blob/go1.27.0/src/internal/filepathlite/path_plan9.go#L33-L35) · [Pinned Windows source](https://github.com/golang/go/blob/go1.27.0/src/internal/filepathlite/path_windows.go#L184-L198) · [Lean model](../../../../lean/GoLean/Path/Filepath/IsAbs.lean)

Pinned native Go on darwin/arm64, the simulated profiles, and compiled Lean agreed on 532 boundaries and 60,000 samples with seeds 17, 2718, and 65537. Native Windows and Plan 9 execution remain unverified locally. [Differential comparison](../../../../difftest/path/filepath/is-abs/check.go)

The native entry point and all three Lean profiles passed fourteen fixed fixtures, including drive-relative names, both UNC separator styles, and Plan 9 devices. Accepting any drive prefix as absolute failed two. Accepting a lone root without a volume failed three. Applying the Unix rule to Plan 9 failed two. [Reviewer control](../../../../difftest/path/filepath/is-abs/review/main.go)

```sh
make verify TARGET=path/filepath
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/path/filepath/is-abs/review/main.go
```
