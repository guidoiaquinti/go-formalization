# Independent review of `encoding/binary.Uvarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `c19343f35b3011c86a2cea1b04c20b2468c782db43e20bb6122fa412c7588623`

I traced pinned Go 1.27.0 `src/encoding/binary/varint.go`. The unsigned decoder accumulates each seven-bit group into a 64-bit word and stops at the first valid terminal byte. A missing terminal byte reports zero and a zero count. An eleventh input byte after ten continuation bytes reports -11. A terminal tenth byte may contain only bit zero. A terminal value of two reports -10. The model preserves these output classes and reported counts.

The author passed 3,330 boundaries and 60,000 samples using seeds 17, 2718 and 65537. Eleven independent native darwin/arm64 and Lean fixtures passed. They cover empty input, a noncanonical terminal group, ignored trailing bytes, truncated ten-byte continuation, eleventh-byte overflow, overflow and success at the tenth terminal, and both signed extrema. Rejecting a tenth continuation byte immediately fails three fixtures. Reporting the whole input length instead of the consumed length fails three successful fixtures.

The scope covers slice decoding and returned words and counts. Streaming readers are separate declarations.
