# Independent review of `path/filepath.globWithLimit`

**Verdict: FAITHFUL**

Evidence SHA-256: `984f07a296dabf2ea900300aa37106842ffcde16443f4482a4db66feb7b8567e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `14f417004d9766ccee45a3909c1da364350a75e99e7ba5ae2d501ea513f4ca77`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `17483f84ea43b03d6d3216042b5b14e031a6e920458ca9ab8402ad44ee32f429`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `508d9697a63e01bc16b1661c449ad62afa0bf94407ba4f45e39295481bad67a0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `28b529bdccc89ebe9a7da878334dadcf5f7307f95f21a4e341fe7d0233906559`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0aada9112704a61bae0f268183377f18d08b24edb86244789edd63ff9d788a88`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d791879a8405c26d9e06bf6abf925400ef660f2c70bfc4a1a234ddb08ec8e741`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `25bb9bbfe8ff62cdb9691c3d651bd465bf598d32236dc75ddef6a9f4b90ccf32`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6445f2866916ddff37311217491b521d14ef3c1558c83439227d9713666d754d`.

I traced Go 1.27.0 `src/path/filepath/match.go`. The private function checks exact depth equality and performs its initial Match validation before examining path metacharacters. A literal path consults Lstat and is returned byte for byte when it exists. A pattern splits into a directory and final component, cleans the directory prefix, and scans directly unless the directory body has metacharacters. In that case it descends with incremented depth and accumulates directory matches in traversal order. The Lean recursion follows those same decisions. A recursive split removes a trailing separator, and the equal-directory guard prevents a stationary descent. The Windows drive-relative cleanup appends a dot only after its volume prefix, leaving no metacharacters in the remaining directory body.

The author checked 271 native boundary patterns and 60,000 samples using seeds 17, 2718 and 65537. Twelve independent literal fixtures passed both the pinned darwin/arm64 private overlay and the Lean executable. Cases distinguish entry depth 10000 from 10001, a direct leaf scan at 9999 from a recursive descent, recursive success from depth 9998, raw literal spelling, absent and nondirectory results, malformed patterns, and accumulation across two matched directories.

The model takes filesystem query outcomes as explicit inputs. The respective native filesystem and syscall behaviors, and native Windows and Plan 9 execution, remain outside the review.
