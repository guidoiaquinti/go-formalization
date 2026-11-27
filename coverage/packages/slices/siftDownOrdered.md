# Independent review of `slices.siftDownOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `42ab7ca668fc318863b5a3f84852dbcfef826f316caf3f3ec90aafb097ab8313`

I checked evidence commit `d508653` against the private Go 1.27.0 helpers and their Lean counterparts. Sift-down first compares the available children and then compares the root with the selected child. The independent fixture compares (2,3) before (0,3), exposing a root-first mutant. The sole-child fixture makes only the root comparison. A panic before the root exchange leaves the backing untouched. The source and model keep changes inside the given subrange. The model carries the typed panic and trace while suppressing all later comparisons and writes.

My `go run ./difftest/slices/local-sorts` replay passed 12,963 boundary rows across the three pairs and 20,000 sampled rows per pair for each seed 17, 2718 and 65537 against the pinned private Go overlay and compiled Lean. The independent `go run ./difftest/slices/local-sorts/review` passed six literal fixtures, two per pair, and checked insertion rollback, sift comparison order and heap drain exchange mutations. [Independent control](../../../difftest/slices/local-sorts/review/main.go)

The native observations cover valid int64 subranges and representable heap indices, plus comparator scripts that may return inconsistent results or panic. Callback mutation, invalid bounds, signed index overflow, concurrent access, and other generic instantiations are outside this replay.
