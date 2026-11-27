# Independent review of `internal/bytealg.LastIndexByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `d16f388568cd5242ebdc31d658000e220ac9cc7dfd8f781f912ee34b52f6fd9f`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b6c0f9d30b950dd2a3db2c5b79e7b74ea7b6df8774baf58dee9bb1fe33b7b4e8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3e82103b87a74a77234f7f7c1acc3f6651b4e71500317b92657707999706355c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c2405563b09212738315cebad2141baf007603717fea11b5245ae8db6e04b55`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `50253168b288ea59cb5b4655e2779996c4af3b9c6a3158d8241a7f8731343b79`.

Go 1.27.0 returns the last matching byte offset within the supplied slice, or minus one if absent. The Go body in lastindexbyte_generic.go scans from len(s)-1 downward and uses the first match in that order. The Lean result is characterized by a valid matching index and the absence of an earlier or later match as appropriate. The missing theorem is equivalent to target exclusion from the entire list.

The author comparison covered 3,587 boundaries and 60,000 samples using seeds 17, 2718 and 65537. The shared independent control passed ten literal cases against the pinned native declaration and Lean. It checks nil and nonnil empty slices, empty strings, repeated bytes, embedded NUL, malformed UTF-8, offsets inside valid UTF-8, misaligned bounded views with adjacent sentinel bytes, and vector boundaries through 257 bytes. Returning the opposite matching index would fail five cases. Returning zero on absence would fail four.

The claim covers the returned offset and absence. Instruction layout and native execution on other architectures remain outside the review.
