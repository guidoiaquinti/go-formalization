# Independent review of `strings.ToLowerSpecial`

**Verdict: FAITHFUL**

Evidence SHA-256: `6f6203ab0e6a3bf56e8afddfc18375bcb0dee40cea980cb9b259fa8cdec3576e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b5e694272c9498e8619d9510d57bd0db7b6d72eaaa47fc593ddef700007a18f0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3457019c4b82d84eb7f51e5e3e4ad3e84c88d322acda9fae06071dc1373ede32`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `589dff8bcf554dad60bf56553069b3016534764b8ab45ba54d0c614e07a3c3b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1165beba94dc7cd5014d814cfd6ccebe0ba0adc7384307eeee54eac9a9210ac1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ade32c1dd43731c4650e72853042f93c4b026d892d86907180368a3e6f72521`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a4fa68486b8fc45637ec1ce79a608b502f30f94d21397314ef66c3620ada2ddd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `19044767eb0eb964a7ee94e9014efd73ae0087c54b92c004d4c3712c7e8be37f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `941eed1646183d7fd3293be9d5af8c230fe5392fdd88dd104da442d2e9a3fd14`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5f0858f100a56911ee3e8f8d0154dd588549e5f6de425d3d6d3aeeb457c95c1d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e405384fbfea77a3a1066ff8658454e74caa4f2f704f90b145c1c445f6e3984a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3142466495bd7f0e0c8a0c498c76d7e4b4f498d5e3d0109aaf4bbf4603b37c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4406bd2709a501af47e9e3253b5174dd6969c19987a2c93f90e7f9d57f7e3abd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `161bb2e0d84dbfb46244f8919e24e8a4be74130d7bf0761eced3ff95c47c6dd4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ecec8cbeb7fc3abf31b4b8d52e735617810c764d55db4a210f93d1d900651afd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8a122163c7b884d5005cf26395533901cb70a757a35e0c1666737c6a75fc49d3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8b96c348b794db9468112f94020e880a10f7c9ec3bc1bc6319946ebe7144ce45`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a4f45abd4578ccbb87f0bb39211c6a8bd01c238c1e0b3f663d8c0b293eceee7d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `65e0525eafabd72bbe9083245e6f60ea954612846cd4749da1476317b067507b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8c6ea9a10bfa9ee4f16f40a22b264720937f0df7e16d4c94fc072301168d1baf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `75c10c5fd2feae07c209b3a76d0bba858ce828496567c73b2c5df49c171bfab0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `937760838b03ad39cd34110b98a45ebaa223fd32dbc0ac8c00f6e943618b935d`.

The review covers pinned Go 1.27.0 output bytes with supplied SpecialCase tables. Reviewed d750343. The claim covers lowercase mapping and raw string output.

The exported wrapper calls Map with SpecialCase.ToLower. The Lean half-open binary search preserves the pinned midpoint and comparison order for overlapping, unsorted, and reversed ranges. A hit selects the lower delta and is terminal, including zero delta. A miss uses the pinned lowercase mapping. Signed 32-bit endpoints, wraparound arithmetic and the alternating range's lowercase low bit correspond to Go's rune arithmetic. Map drops negative results, encodes invalid nonnegative results as replacement runes and traverses malformed input at one byte per invalid unit.

The author differential has 1,232 boundaries and 20,000 samples for each of seeds 17, 2718 and 65537. An independent fixed-output control checked eight pinned public Go and Lean observations covering midpoint overlap, search direction through an unsorted table, fallback suppression on zero delta, signed range endpoints, negative drop, surrogate output and malformed decoding beside a literal replacement rune. The overlap row separates the expected result from a first-match traversal and zero-delta fallback. The unsorted miss separates binary traversal from scanning every entry.

The scope excludes allocation behavior and unsafe backing-storage observations. The replay is finite; the Lean theorems establish selected model equations.

The registered control uses the repository’s Lean driver lookup and runs without command-line arguments.
