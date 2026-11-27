# Independent review of `internal/bytealg.LastIndexByteString`

**Verdict: FAITHFUL**

Evidence SHA-256: `76bab6077e0c959737055bca016eb6771e210646e1c0f7288251e5e673195349`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `29acc895dc1ea45290d0b6a16bc41a97962c9e938fb28ed4ce1d89c3a04ccaf1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0272f47bfef81bb5323deb9c1ceb25167c10f23dedd2c74a37deeab6e7248629`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `67f62f9f26be2c643729ee206f90aa3059d0e7eddaa014aacee366ad805a38c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6fd3591d0a23ef4f8236feb784bd8962fa4d154154032a6d1d249f24a3627402`.

Go 1.27.0 returns the last matching byte offset within the supplied string, or minus one if absent. The Go body in lastindexbyte_generic.go scans from len(s)-1 downward and indexes raw bytes, including those inside or outside valid UTF-8 sequences. The Lean result is characterized by a valid matching index and the absence of an earlier or later match as appropriate. The missing theorem is equivalent to target exclusion from the entire list.

The author comparison covered 3,587 boundaries and 60,000 samples using seeds 17, 2718 and 65537. The shared independent control passed ten literal cases against the pinned native declaration and Lean. It checks nil and nonnil empty slices, empty strings, repeated bytes, embedded NUL, malformed UTF-8, offsets inside valid UTF-8, misaligned bounded views with adjacent sentinel bytes, and vector boundaries through 257 bytes. Returning the opposite matching index would fail five cases. Returning zero on absence would fail four.

The claim covers the returned offset and absence. Instruction layout and native execution on other architectures remain outside the review.
