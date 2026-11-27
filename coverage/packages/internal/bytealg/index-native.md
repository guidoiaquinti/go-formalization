# Independent review of `internal/bytealg.Index`

**Verdict: FAITHFUL**

Evidence SHA-256: `9459fe6e4193466af5df891cadb70c8036db815dff68b790dfc1640d14f6a8c2`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `88c71c9088e23b14d78474fe251d1b639825a791634c16b064d4d7d4f37187e8`.

The comparison calls the pinned darwin/arm64 `Index` assembly entry through the same-package bridge, using legal small needles. The source constrains arm64 to widths 2 through 32, and the checker uses widths 2 through 31 on amd64, the minimum of its dispatched maxima. The model adapts the existing byte-window scanner and proves the first-match predicate. Returned offsets count bytes in the supplied slice.

The author comparison checked 7,936 arm64 boundaries and 60,000 seeded samples using 17, 2718 and 65537. Eight independent literal controls passed the pinned native entry and Lean. They exercise overlapping matches, malformed and continuation bytes, final and equal-length matches at the maximum selected needle width, independently misaligned views, and matches beyond the end of the supplied view. Returning the last match, dropping the last visible haystack byte, or extending the bounds to include the appended sentinel match would fail two, three, and two cases, respectively.

Return semantics cover the stated caller-supported domain. Architecture-specific instruction layout, wider amd64 AVX2 widths, and native execution on other profiles remain outside this review.
