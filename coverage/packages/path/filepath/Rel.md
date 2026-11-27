# Independent review of `path/filepath.Rel`

**Verdict: FAITHFUL**

Evidence SHA-256: `1ec63d83c73ff92b880fbce10aeea8b10b27fe7371e522817d8ade91011552ac`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `97deb950b3cf251640462febd80d0746b99b43b9cdfd91d38171b39a8b986162`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7bf9c6035a8d1229c9315bf9b960ea370e01ee7ffe57a36a4e06f79416d06a17`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dbb169e1ddb6bd2d5a839fe3c4ffcdd4558b423f076db39d0b27ef51c105eba8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b75688eea2fd176b069d459dadde353966c4a3307a0f3602e7e161f80aae30fb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8bcb04f0591ed906da1bf8583616abed21d5f159b0a2133916b2cde4ded6ed79`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `299511fa9678567b690ea31cbfd5b21a53a6ad0d8df125dbd8f12529fe7a90d4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f7c1a85e9c86ebf9d3626bcefdaeb870166bcc63b33f075e01272b4e1dc9ffb0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `917c4cae911ad63a65bc6f36d6618b0443e903c1ac2acd8c866f0322d6d032a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9dbecf1bc7cfe93b6aecd3e70e1fd48db07d32d52db3b5a75710c794e58615c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4240035f7fc73b16545ee2d227e87cf3ffa40a237c0809ed3e89b07270cc865c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0ef2e5078a4c23cc22b2fb16738323b316539ff1d5010200bbd40e319fc520ff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9b9246e36ee09a010ec0e62898a03b849789dd856150d67639707585e52ace3c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1140ed290b79425484a61fddf98c8583e77bf77cd457c747431320faa1c91e81`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7397a67a010248e7dd89f2626c5b3f044e4e7d833d58331d0f4243fdd05afd78`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6773ccfab2325663c9d609098a326ee3b66189c23554f9f0f31eb4494efc12ff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `94beeb118e7eeaa490f5eb934aae9ac7b3b7e4abff986c14d7f66566bbddc79c`.

The model matches the pinned lexical source. It computes original volume names and cleans both paths before comparing platform words, strips volumes by byte count, treats a bare UNC base as rooted, and rejects root-parity and volume mismatches. It skips common components, rejects an unmatched base `..` component, and cleans the synthesized relative path. Errors preserve original operand bytes, including malformed UTF-8 and NUL. [Model](../../../../lean/GoLean/Path/Filepath/Rel.lean) · [Comparison](../../../../difftest/path/filepath/rel/check.go)

For a bare UNC base and the same UNC target with one trailing separator, the whole cleaned paths differ, then operand preparation produces two separator-only bodies. The component cursors consume their separators and become stationary at the end. The original extracted Windows code did not return within a bounded subprocess observation. The independent cursor trace shows the exact non-advancing state and Lean emits `diverges`. The reverse direction returns the original root-parity error, and identical paths return a dot before component scanning. [Finding](../../../../FINDINGS.md) · [Reviewer control](../../../../difftest/path/filepath/rel/review/main.go)

Pinned darwin/arm64 matched 529 boundaries and 60,000 seeded cases using 17, 2718 and 65537. The portable extracted pinned Windows scan matched 60,004 inputs, including 10 stationary states. The independent control passed 12 native Go fixtures and 26 fixed Lean results plus the separate divergence observation. Cleaned error operands, Unix case folding and omitted root parity disagreed on two, two and one fixtures respectively. The native Windows subprocess control cross-compiled, and native Windows execution remains open.
