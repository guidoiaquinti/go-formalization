# Independent review of `path/filepath.FromSlash`

**Verdict: FAITHFUL**

Evidence SHA-256: `3efbbb9cc9dd841eac6eb5b757095e186861fa85dc58042647df14fcc5121efa`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1c6a7e188e152ecc424d6be6be96d66d829ef326fcb5cf7d2a660fd73b9c1454`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0cf13c1883877a5d6d82ebaec784d245a31214c5deffb77849e1441779ebacaf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `409a9947ab9db28733926bcbe46aebcd1932d345f79e788693450730c9f38c70`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c5dddfb47b78045d1ddb7c63f2a4a61e45f87727d53d8d7dde9db5bb8905ed5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74cef687d3c8488676f07c7712942686348cd06dc6877f37d033f366a6427ebd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90eba35d7249769cfb021198dfb973739d871af6443f196b724935c952728842`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2577210d5d7d70636446ccac8f0edff528342a72c27955f295c6c4a67d86c1cb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4934e00c2424ea7103ebcc0d9e1eb404a0bde21415d57ec84ff49fc83f7fdb83`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `487a7971503d3049ede4f0c673b3018afc891b95dcda668bf3fd0a7f990881ec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `16ea22dd703da72b958d34e12386a910378750e8d0c2896221daa2ca21753c1e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `25cce4f7681afed8c36cceb1fe9b46c872fc03d3e07da076998a34a996026d8d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9fdff98d788e41f7399d923fee00d968116314c9fa34ee3647211b17632f80a3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `638e1ccc283fd7f50562414d2867550db1f349e5a294a14663ab20ed0c5ae206`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `943ac359ede5fdad12cb0bb9d3e88042dbdc732a87769caa222f290a630c4f8f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98f0bd55f53bd418f6f61234d4c7bc4ab7ebb9e508d343398f1bfe5ba865cbec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a3e69e380fc8f617a1dcb260469e72ca7ccb279e58d22795e88414f68a48d203`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6d0022b02bd68a8579aacaa82ee8b1c7fd32f0c35aea3f343270494771ec9d0d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d61f5857942f61553aa3f695a9bf4b1172ee3f93dc662c426a14b2b6d88a7fb5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e0fb86d0430ad2e5cc5419628a078cabb7b8ef34efa7e56d630d7319533ea23d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76690d11fc8f664834c48cd81c84b69a809769085c1713a8ccfe9c53d82fe663`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cccd6b3cda335ef6a2ec0b268bd59aff0954fe09a3d7074194fc347f30a3b905`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6d979ebd8d15af4345b6361cfc3ff456d8584b9d22bcb81a09b0908b0f20aea5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `409c2c5d24c2813efe650447edc1ef1aa1c43fc3b13cdc1ea4e1c905ce539643`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `112bffd6de1bf35a31aa0d1829d5f386faabbffbcc90b0bff93b5f90a8f87f35`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7c3549aa569296780e01afefb1c9d47ac30f0269289448eeff75efe52e9aa4f7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3b5aa9bf7e743feb1f9db402be07e4db6b9635e5c0db5b344b705e4cd341c5b`.

The pinned wrapper delegates to `filepathlite.FromSlash`. The slash profile returns the original byte value. The backslash profile replaces each forward slash byte while preserving every other byte and the original length. Lean proves the identity, per-byte transition, and length preservation. String backing identity and allocation are outside scope. [Pinned wrapper](https://github.com/golang/go/blob/go1.27.0/src/path/filepath/path.go#L102-L104) · [Pinned implementation](https://github.com/golang/go/blob/go1.27.0/src/internal/filepathlite/path.go#L184-L204) · [Lean model](../../../../lean/GoLean/Path/Filepath/FromSlash.lean)

The pinned comparison passed 520 boundaries and 60,000 samples with seeds 17, 2718, and 65537. The public entry point ran on darwin/arm64. The Windows separator profile was exercised through the compiled Lean model and a Go replacement oracle. Native Windows execution remains unverified locally. [Differential comparison](../../../../difftest/path/filepath/from-slash/check.go)

Seven fixed fixtures passed through the native Go entry point and both Lean separator profiles. The cases include repeated slashes, existing backslashes, NUL, and malformed UTF-8. Keeping forward slashes unchanged under the backslash profile failed four cases. Coalescing repeated forward slashes failed two. [Reviewer control](../../../../difftest/path/filepath/from-slash/review/main.go)

```sh
make verify TARGET=path/filepath
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/path/filepath/from-slash/review/main.go
```
