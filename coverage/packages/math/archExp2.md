# Independent review of `math.archExp2`

**Verdict: FAITHFUL**

Evidence SHA-256: `9c23c2ac3ddeb895db83860673128ffffc89533fcfb1671fa8f49b3915f73a7c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `68abe5741c4de0b6e45a06a1683043a2b469d63e9acee69ffeb50031f68074c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e1b3d1f0f68fcbb42f1bd9d71a141714b6255c42787c6fe35e7daede554781a0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5603972da8fd79a6313ebfe0762288882d5254f4f632a35653daa49b10ceca13`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2651912fcab892e375dc9ee4bf365eebdff93b7ac5f53e542a6a49cf4589ee43`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dcc9d104606b10be6711a3531af598a37bf238bf02991e86d4f8f54525fad055`.

I checked the pinned Go 1.27.0 ARM assembly and the fallback build
constraint. The native ARM cutoff returns zero immediately below -1074.
The fallback body instead has a false availability flag and panics with
the exact string `not implemented`.

The independent control checked 34 Lean rows in the ARM and fallback
profiles and 17 corresponding native darwin/arm64 calls. Fixtures used
signed zero, infinities, a NaN payload, integer powers and neighbors at
the early cutoff. The author's native comparison passed 2,170 boundaries
and 60,000 seeded cases with 17, 2718 and 65537. Its independent reference
differed by at most one binary64 step, and the runner counted 243 early
zero results.

Native fallback execution and the LoongArch and RISC-V bodies remain
unverified locally.
