# Independent review of `path/filepath.cleanGlobPath`

**Verdict: FAITHFUL**

Evidence SHA-256: `af2afdbaed7324e06084a2a91de5c195c8507abd87105aec08aa1e88add7d5f2`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `03fafc315b62af6a8af2767a344feef98f0bc73a74dd1d0ff6c4d8207fdff7f3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `28c6b50509287a36187c23ed477a240e03267c010700a83f4960521740538b4f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6907248f2401742a7eebe524cbb04908ecb52b95a1860535ebf5caf30dd577fa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `17c22f32d21c9a61b3794db34d51c0cfbb15f79abfca9b044c65ff1e5e7e409d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df7feaba565afbdb3823461d2db1a43155d497bb8d37ffc8a4cc3281d27a2df8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `83bc410f4f283bfac8357cfbd3248ea5ff5105789a9e0b412da1f05bf415cc33`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01e8b7d9500d9f10ae8a2a05e0038e1b9b18cbd2bd1f4b190bee9b754d217803`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3cc505dc35355cf237a22f792b32578387c1043348b5e84dde1c546412b6e44a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `909c62a26cfcac111ed08092952cedb9161f61c6a415aa20b161e3e9cd15fc52`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7290cc43000dd225ca320c72086765dc6d9d948a7825d2040f1cef6583e6dbe7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5291fd1da0a8019e50d36b99004c7c9ea674fe2c25e65b5e7584a99e1d551165`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4424d83af87db86a1b3ce12f89254aaa1d3a30eb34ce78962eb4378e370f7fea`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2a6cf74b9fe63b239a73989e1e7c1580eb4c4d653d3f4ecee3776252bcf6fc31`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `97ba341c728f8a8f87357928a2197c380021297c7475a47407d3b0c7b3ac0fef`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0eeae8fd3d5244fdafd303424cf7a22bdd8a8d507aab5fb8b2f804fef365f021`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5827a2036b7c8da2f9329a82697b2892c93562c1b5cb98c11f9b3824cac4e94a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b822b810f82edf0d7b8533638170171a6223374467b969ade46842ae412401be`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `38458295d3bae726d4e56dd425d3dbc0f5401d765a09f0bb54560f9d9b17eee6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a97fadbd686097a3d0b4926252eb382af3544d95a76747cd74740bb9c5cb48cb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bd3b9baddfc648b8961ae721de56a5624e530a59aaa5c41669fc007e38f5b3d4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6ce9bd5ea7cd3777d15c0018b818417801949fcedb42adff41eac25026bb7946`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98d64c01709dc85f538d3f9f76bba9cab167fe1f12fda7c403805d31cbf25635`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2474535b77a27c7ffa9a36589c9afb882fba320e822c04e5a4a7f18de10413f7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `86c76d62aef6222bfc3b8df5ddd6d311d8833fc304f00e18edb4449c16c8dab1`.

The pinned helper returns a dot for empty input and retains a lone native separator. Every other input loses exactly one final byte, including a final byte that is neither a separator nor an entire UTF-8 rune. Lean proves each length case. The comparison exercises the host separator. The separate Windows glob-cleaning declaration is outside scope. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/path/filepath/match.go#L285-L295) · [Lean model](../../../../lean/GoLean/Path/Filepath/CleanGlobPath.lean)

Pinned Go through the private package overlay and compiled Lean agreed on 520 boundaries and 60,000 seeded rows with seeds 17, 2718, and 65537. [Differential comparison](../../../../difftest/path/filepath/clean-glob-path/check.go)

Pinned Go passed all nine reviewer fixtures. Requiring an actual final separator failed four. Removing the last decoded rune failed one where trimming only the final byte of a multibyte rune preserves an incomplete encoding. Fixtures also distinguish the native separator and its opposite as one-byte inputs. [Reviewer control](../../../../difftest/path/filepath/clean-glob-path/review/main.go)

```sh
make verify TARGET=path/filepath
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/path/filepath/clean-glob-path/review/main.go
```
