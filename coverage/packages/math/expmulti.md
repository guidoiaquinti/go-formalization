# Independent review of `math.expmulti`

**Verdict: FAITHFUL**

Evidence SHA-256: `46a2dc4aecfff7e3fcd3d660cf58613060c78552dc2155e82f82cf1743d31937`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fb67b97132f12e47c256547b21411b29c3cc0839186c245cdb8cac54b1a992e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `50ea7e64708b0ffa83436f59af97934a7844ad4b1f139c9b8e0c4b6f5ad6e369`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1c8a4ba9f3d39b5523a3eccc4fba226c0db9cdabf5152fbff83b03d0c0614ccc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9315edcf64e1d0b9abeba592a1193f767d8a42821cd5978e5afe6771a9c4f9b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b437c455e7c2f32b9370d70632f12fd231e7bdbee0d1c3509d8814e74cefc84c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad1c957b7cf2376b46cdc57c45b61f91b3265ab2f3399e372a3bf4455bfeab8e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ce86f91f330f40fdeedd4c4f0d44a9116f55629ad4ff9c397c5695802943532e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c0c9056ddb374eb02d26aa44ae9dfb6ba7bc6e4c70784838a90152c82a98c4c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `42192915b3fbb5d4d61653bb0d569140b0fbb44149d59aad65442a58bb858ab9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e797a5b9592dc66486f5e5f9ada21b44db08176f9872c3e9ce702dc8d6df354a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4a22fa498011007d2cbf1ec659b2c04803663370715852b02f606398486321f3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8e9dd2b073b8eee71bd68ae2c4c51e67af5231a70fdc7eaf797f547b78bcf890`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6da48a20fefbb2a5f0f7a041abf33328c9a2d5e73df06a66fe0d52087d721646`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `abcfc2bddf8aa2848fbb7d0bee3ed95bfff5b2a34f7cac039a953b5867efa198`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7d61d2e1df146f2cba116259e820672d07b5f980d1c87a078bee3a006d671a9e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `765a1ae4ef92a23c37c69a5bfbb1f05a04c91ca13ca0de147b34397b519b27e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c7b2a1abfbae44401daa6f878c626846a049407881b660a5d95c6cb6e68bc0b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `485febb88d24c95b2c1f34c234a04f36ee6c686413f4996d601939a498b67a5a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d492d3920cd06ff8bd184ebe5ca7dd50451ab6240358cde4c89cd15bc98d2a0f`.

The model follows the pinned helper's subtraction, squaring, Horner polynomial, rational correction, and scaling of the rounded positive intermediate. Go 1.27.0 compiler output confirms split arithmetic on amd64 v1, four fused Horner stages on amd64 v3, and those stages plus a fused subtraction on arm64. All three builds use the generic implementation through public Ldexp. The model preserves its signed exponent, subnormal rounding and overflow behavior. [Model](../../../lean/GoLean/Math/Expmulti.lean) · [Scaling model](../../../lean/GoLean/Math/Ldexp.lean)

The direct private overlay on darwin/arm64 matched 100 boundary rows and 60,000 samples from seeds 17, 2718 and 65537. Its largest difference from the 256-bit reference was one binary64 step. Lean proves exact scaling for the zero reduction and exact words at the zero, smallest subnormal and infinity boundaries. [Comparison](../../../difftest/math/expmulti/check.go)

The independent control exercised all 2,100 exponents at zero reduction and 162 nonzero boundary combinations. All 2,262 native observations matched Lean, and 6,786 Lean observations across three profiles agreed with public Exp and binary scaling within three binary64 steps. Incorrect variants that lose the exponent sign, ignore the low argument or flush subnormals failed 1,165, 78 and 104 rows respectively. Native amd64 execution, nonfinite or larger reduction arguments, and native int extrema are outside this review. [Reviewer control](../../../difftest/math/expmulti/review/main.go)
