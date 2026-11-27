# Independent review of `compress/flate.histogramSplit`

**Verdict: FAITHFUL**

Evidence SHA-256: `712c4370d69210ae30c82ae168f230d34deaf40a583753672ec1fa586cd0dc03`

I reviewed evidence commit `7baa252` against pinned Go 1.27.0 `src/compress/flate/huffman_code.go`. Pinned source reslices the destination through capacity, consumes the remainder and then steps through four quarters. The four updates remain independent when several lanes select the same bin. The independent 8,195-byte input detects lost remainder and lost repeated increments. Traversal order is checked in source and the final frequency result is checked by replay.

The pinned same-package Go and Lean comparison at `go run ./difftest/compress/flate/histogram` passed 443 fixed rows and 20,000 rows per seed using 17, 2718 and 65537 on darwin/arm64. Five independent rows passed with `go run ./difftest/compress/flate/histogram/review` and rejected seven incorrect observations per declaration. They verify that capacity 255 panics before writes even for empty input, that a zero-length destination with sufficient capacity works, and that counters wrap modulo 65,536 while the backing past bin 255 is preserved. The isolated kernel audit passed. Scope covers disjoint source and destination and legal finite slice geometry. Panic payload and concurrent mutation are outside scope. [Independent control](../../../../difftest/compress/flate/histogram/review/main.go)
