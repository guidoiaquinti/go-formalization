# Independent review of `internal/bytealg.IndexByteString`

**Verdict: FAITHFUL**

Evidence SHA-256: `7fdbc5db64bcf1c09da1acde1e9debb19e874b0808dcfc8df656eab88b4fbe37`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6d5c0915a1f57d4c1ed80897e865418d71cb1b6f38b8c68532c2ffe772ce30bd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7329f718b373de02bfdce4e67d06ec0d79b9272827c91c42b35a61fd95dd3ede`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `92be87d24150bd98bd07ecd04bef036d502db465090317dffe0831e5788c9d4f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d0f95e9b0b57708b8887c8ddc4e96c57ccb9f218c2e4ff6c4c5489d684678574`.

Go 1.27.0 returns the first matching byte offset within the supplied string, or minus one if absent. I inspected the native declaration and arm64 entry in indexbyte_arm64.s. The compiled darwin/arm64 overlay contains a direct call to that assembly symbol. The reported index counts bytes even inside UTF-8 sequences. The Lean result is characterized by a valid matching index and the absence of an earlier or later match as appropriate. The missing theorem is equivalent to target exclusion from the entire list.

The author comparison covered 3,587 boundaries and 60,000 samples using seeds 17, 2718 and 65537. The shared independent control passed ten literal cases against the pinned native declaration and Lean. It checks nil and nonnil empty slices, empty strings, repeated bytes, embedded NUL, malformed UTF-8, offsets inside valid UTF-8, misaligned bounded views with adjacent sentinel bytes, and vector boundaries through 257 bytes. Returning the opposite matching index would fail five cases. Returning zero on absence would fail four.

The claim covers the returned offset and absence. Instruction layout and native execution on other architectures remain outside the review.
