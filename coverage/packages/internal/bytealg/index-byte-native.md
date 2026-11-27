# Independent review of `internal/bytealg.IndexByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `4aee35e7039668fa5ad6d4b873080e574f4d2df0dd21d61ff6e0f073a1e611e2`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9203387645ca726c391cfcec59fe43f6060bcf567fde27a8549d0b0ef917b98a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fdc2b4c5afc44d805711310d405bf5f24599568472fe928fae7486a05d88aa45`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c1cefe846ea0f4c9297396867d72039b99d633ffb84f9a19347be5e874a3811a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0cbeb96a1bc3e404eb9bee85946faacf1837224166beffea76c5a22991342a21`.

Go 1.27.0 returns the first matching byte offset within the supplied slice, or minus one if absent. I inspected the native declaration and arm64 entry in indexbyte_arm64.s. The compiled darwin/arm64 overlay contains a direct call to that assembly symbol. The Lean result is characterized by a valid matching index and the absence of an earlier or later match as appropriate. The missing theorem is equivalent to target exclusion from the entire list.

The author comparison covered 3,587 boundaries and 60,000 samples using seeds 17, 2718 and 65537. The shared independent control passed ten literal cases against the pinned native declaration and Lean. It checks nil and nonnil empty slices, empty strings, repeated bytes, embedded NUL, malformed UTF-8, offsets inside valid UTF-8, misaligned bounded views with adjacent sentinel bytes, and vector boundaries through 257 bytes. Returning the opposite matching index would fail five cases. Returning zero on absence would fail four.

The claim covers the returned offset and absence. Instruction layout and native execution on other architectures remain outside the review.
