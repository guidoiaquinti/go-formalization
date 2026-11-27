# Independent review of `path/filepath.abs`

**Verdict: FAITHFUL**

Evidence SHA-256: `835c05af225318e688f0f53bc0293a64f17e114adaaedb80de8fe6f633839f4c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `753e393fd9944b65371bcec1baadb2348b3c3ef1517ba09a9a255d2e17b61837`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8454f5e77ef16148126ff948e7fc73cb9ad06607796aa28cc019a300bb26e3ac`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0e5e1ab011d631925fa417bf6f6d8d6d2d08450d15e71ef704e7fcb5fd7039df`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9d8c9173d519c3d6666c863b8aa994e875b358d0ce28b9d54a50fae2c263b797`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aa547e8fd4526bc530596c1e1e09a70f9870e93cc1a920742fd9dbfb8efc4278`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `522610fb48e7aeda9d0e93db1272de213723e6800074192b880dc91269e49911`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `be431a2c551c116e9722ac3119ce82d62bbbf477794bf47d90c88e2652486e40`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `52c22166b09560691df537aa26b10932d853cd50c989eb37efc8aa3fc6588978`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ebd5c2baaa1a3519f52b8267f02fe8524f5e526f6d1553e23c7c0db859d537fd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `05026d589780d0b7f88830c6554236fa9fd67d04f90864f89c74d78172a26f79`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ced1ee4ff82a53cb81747f154665e422bb6c89abd92a66b8ecf3cfdcf492c1e2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7fa0ae2996c905d5e559149f4b9a2a75967c0ab248a613bede5bfef189f1a817`.

Reviewed 5ab5b78 against pinned Go 1.27.0. The Unix and Plan 9 implementations dispatch to unixAbs, which cleans an already absolute input before looking up a working directory. The relative branch propagates Getwd errors and otherwise joins with that directory. The Windows implementation calls FullPath on the original input, replacing only empty input with a dot. It preserves the FullPath error and cleans a successful result.

The independent control checks five actual native calls in a real temporary working directory and ten model profile observations. The model checks Unix and Plan 9 working-directory failure for relative paths and bypass for absolute names. In the Windows rows the supplied FullPath accepts only the exact expected argument. Cleaning `a/../b` first, passing an empty argument as-is, or bypassing FullPath for an absolute path would each disagree with a fixed observation. The control also supplies a raw NUL and distinguishes successful-result cleanup from error forwarding.

The author comparison checks 522 native boundaries, 60,000 seeded inputs with 17, 2718 and 65537, and ten fixed profile cases. The reviewed scope covers dispatch and the supplied syscall outcomes. Native Windows and Plan 9 execution remain unverified locally, and the Windows driver cross-compiled.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/abs-private/review/main.go
```
