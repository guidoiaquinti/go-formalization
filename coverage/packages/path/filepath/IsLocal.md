# Independent review of `path/filepath.IsLocal`

**Verdict: FAITHFUL**

Evidence SHA-256: `83b203d964dafcb4fba907824a671dc57c22d2e28b891aaebbc953fc4434f8bd`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d1a993a7c7855df108ac41965883d8b8a85a0fd946b2b4d4e1c48541dfc3bee2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7be6b086a6d3644929505b451a41ae7c93054f75497d93f7682c2e0b041d0f3b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c008c7db1b030d8c4d75002ffc9ef5ad5e62c4c9b2d278444db97fcbb606e112`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30661142bdcb5df8f990a0f3a0668ee00cf84f6552584d6fece5818e54014daf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `adc2a5d158dc794223f43ad53d66d992c12464118c58d97c4644abb5ce5b8ce4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `758e62e534c680d9d68972989e774468837cbdea1e0aacc5e7d76e1e4d3bcfd0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9404dfe8812ef849147c33c6966de6b537fd8131c11ba0afacd25410de1163ff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9ff1d906f2f50c6c0b53f793c014de121fcfdfceedacab64f50fb4037a147d42`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c9200296e67175bef53fe674e156d0a4f84a64e3a5510d7809d0d0cb0e510a40`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76a396335adc989597abf3cd22076f69bc2b6e9433793a22235409a628f5580b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `46b709dc1c940eb0fe289c82273b5ef83147cf547e31bf9172bffbb8eb5e2a2e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `598a73a03323d1e65f01640457cabf44d573a890537fb8ed7f8f8483025bf17f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `726d8e38406ba5f021359638dd808f9d5200767f617781c6c01e8d2ba9db9667`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `83e672f063a15a3f620587c596b800cd2e924ac2ae9956944be218f0c291e56c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `787a92ca0e4046ce513d56b7d45ae92c0c53bf796b6569b0a0520a270e937e8a`.

The pinned Unix and Plan 9 paths reject empty or absolute input, clean paths only when they contain a dot component, and reject leading parent traversal after cleanup. The Plan 9 absolute check includes a leading number sign, and the Unix check uses a leading slash. The Windows branch rejects a leading separator or any colon, checks lexical device names before cleanup, and applies the Windows separator and parent rules. The model distinguishes literal reserved bases from names that require an OS decision, using the original raw component as the callback argument and excluding embedded NUL before that call. [Model](../../../../lean/GoLean/Path/Filepath/IsLocal.lean)

Pinned darwin/arm64 matched 533 boundaries and 60,000 sampled paths using seeds 17, 2718 and 65537. The author also verified 63 fixed Windows and Plan 9 observations and cross-built the Windows checker. The independent control matched 8 native cases and 29 Lean observations across the three profiles. It distinguished ignoring the device predicate, altering its raw argument, invoking it on NUL-bearing input, treating all superscript digits alike and using the same absolute-path rule for Plan 9 and Unix. Native Windows and Plan 9 execution and the OS device-name predicate itself remain outside the local evidence. [Comparison](../../../../difftest/path/filepath/is-local/check.go) · [Reviewer control](../../../../difftest/path/filepath/is-local/review/main.go)
