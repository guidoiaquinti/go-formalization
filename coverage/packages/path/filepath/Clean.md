# Independent review of `path/filepath.Clean`

**Verdict: FAITHFUL**

Evidence SHA-256: `4907b2dd4f3d06be17facc3affa6c0172b16bbe1b722750d70db3f0409960745`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c8d884d386bfa5534476939f1ab32ca3b102736f9cd93545e018320d19c9650f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b34487b7105500bf2d43177e790c093b048a90a05da07a400a8c07558a97893f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `77d94c2d7f4b87bdcd177b75601bf911e54cc9da8fe9e1e334747e73f6b603c6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cc30f340a23f7af737c13eb53ceeeedbbc29e7523398d98241741137323dfd56`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2bd73cd2d0c50164823885c5edcacd0edf8fe574ac54c12ab1d99debaf23c648`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e213621e6caf24b0a5d848cf6d539ea4d82e284b987ec75b60af9d559eae547f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7f8fa904f65a2203a7ab49c2668325fcbb2d5f8ca87dd923b23aa7c9f088f09a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `27ed0bb6cdc3df68bb27a36a3ea08a9ebbbad94f54ec97e1651f788732783acd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2104d38e1bd9756e209fcffddfc18b9343671138298ec13213d1148b63bcc960`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `94e0e9a29c02c101a9663c86a7914a3116c52008eb364ca27b1c249f07a4392b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `312fd0c107303e2722f1c286e31d2a82cc6980349101d1d7473cf34e0ef666d8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `46f9dae6875d84e18838b45565f18c21ad18835611f33226516dd8b4249d64a0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef334635c70865c2729771472a002d49680272b5b15af44caf85e8fd9c301c24`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e40ba43dbdc13ceb8af314eca96b9bece6cb26033d460172ab390101385bdf45`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `503b99eb24fcb14df2f4d7335f783ad00a2dbad127ff203b0186a212e9f7f060`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `66697f5c9217143fd039fbe10bee6cf557dc082acfe410659431daebb9f57b17`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f5324497d96fb15a70fda953d222a69da8da53e7d3b60a717d4ee24ccdba72ee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `184b0241928b41396882fa13f573e5ac44b33ecec3e0d4ba695d492871ae41ff`.

Reviewed corrected evidence 1354d68 against Go 1.27.0. The model preserves the lexical stack and the raw path volume. A parent component rewinds the visible output and preserves any allocated backing bytes. The pinned Windows postClean examines the full allocated backing, so stale bytes beyond the visible output can cause a prefix guard. An unchanged volume-less `foo:bar` keeps its original result because no backing is allocated. The colon in `./foo:bar/..` remains in the backing when the output is reduced to a dot, and the result becomes `.\.`. A discarded question-mark pair can likewise guard a rooted result.

The independent runnable control checked 14 native results and 28 Lean profile rows. A visible-prefix-only guard disagrees with three explicit Windows witnesses. I also extracted the pinned Clean and lazybuf code with the pinned Windows volume and postClean helpers, with portable equivalents of the separator and slash-replacement functions, and compared seven fixed cases and 60,000 generated Windows-profile paths. The earlier model disagreed on 43 of these inputs. The corrected model matched all 60,007. Seeds were 17, 2718 and 65537. The author comparison also checks 532 native boundaries and 60,000 native samples plus 24 fixed Windows-profile cases.

The claim concerns exact lexical output bytes. Native Windows execution, filesystem resolution, allocation behavior, and backing-storage identity remain outside the verified scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/clean/review/main.go
```
