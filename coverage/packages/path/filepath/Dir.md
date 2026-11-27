# Independent review of `path/filepath.Dir`

**Verdict: FAITHFUL**

Evidence SHA-256: `55f6dd8c97cad04026a6b979f7c81413e9ed7e4ddfbee418f14c21e9289adb84`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2695be36340b21ec0df5ea5720463b380e994b90a554c81aa28d9be84a68e4e6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7fe9f9dce073ff019c5341cb99adad519f113b98f3b7a4c8091898e1dca7be9c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `958b7cad8d65f4e8783413783229d6bfb4b9065e06f69b802274ce8e558c1dec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b1d98fa3948d7226bb420da28439e9b333a05652c2ff8de5ef6f5841478a4e7a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90cb873904f24c0288ee60109a0b580684345a28ca4c82765877e7479eddf7ca`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f3cca2aba39d2aa00fe1c2ad2f4a6af8dff01ba8e94e0ac7a941e832568589f6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a61e7f943ec705496be719ab235f7f9eba329b166d2593de1d3519a915bd4188`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `353907619645ba5fc3775abf309f301db129cfd7c4b9b8c3ccb52db6c36a096a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `474b203277b4325cee55d435afe56bc5fb151cce300af6bd80c508df2ba7a537`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a0fe8638b594c52dd02eecd1b523410fc2c568c90c761a8551cb0179932991d9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ba51502695585610540484c455dfc7477f3d8a31b0d2276a291979e4fac893fb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0c4b4b8f944092ac439b59fc0f6c2ff2aedd7c4e8ff31e359d9bb4aac263d4dc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4d4975cbca069004097108e274c7feda11773013fb498fade95bb515454aec76`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0fd56656b01ed429d92cb78541e89cb626029a741c3c45418836bd7857dee469`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7aca9d9ab32e5327d1e01226f4c97dd0295286d43956c2bd90319b19ef3a8889`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `83e2d4c466923447bf223caff6d5dde132d9ff1c31702813022941a44e7948d6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7c02669204d27cba6effa0e3cd254ab879e0813cd6b12ea463cd3aaf842dd306`.

Reviewed df1f776 against Go 1.27.0. The pinned implementation obtains the volume, scans backward over original name bytes to find the directory boundary, and runs lexical Clean on the directory bytes after the volume. The normalized volume prefix is combined with the cleaned result, except that a dot-only remainder following a volume longer than two bytes is omitted. Splitting before normalization preserves the distinction between `a/../leaf` and `a/../leaf/`. The Windows profile inherits the pinned materialized-backing behavior of Clean, including stale colon or question-mark guards after parent traversal.

The runnable control checks 14 pinned native results and 28 model rows across Windows and non-Windows. It includes drive roots, UNC volumes without a directory suffix, exact malformed bytes, and `./foo:bar/../leaf`, which yields `.\.` in the Windows profile. The control distinguishes normalizing the full input before selecting its directory and retaining an extra dot on a UNC volume.

I also compared the Windows-profile model to mechanically extracted Go 1.27.0 Dir, VolumeName, Clean, lazybuf and Windows volume and guard helpers using portable separator and slash functions. Nine fixed Windows-profile paths and 60,000 generated paths agreed under seeds 17, 2718 and 65537. The author differential separately checks 529 native boundary cases and 60,000 samples plus 17 fixed Windows model cases. Native Windows execution remains outside the verified scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/dir/review/main.go
```
