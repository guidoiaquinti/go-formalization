# Independent review of `slices.heapSortOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `121a3d3bd3a6811abab5e042bb16f39716552288621e5774ecfaff8053b97188`

I checked evidence commit `d508653` against the private Go 1.27.0 helpers and their Lean counterparts. Heap sort builds from the last parent toward zero, then exchanges the root with each descending endpoint before sifting the shortened heap. The independent panic after the first drain exchange preserves [91,3,1,4,92], exposing a missing or delayed exchange. Empty ranges make no comparator call. The source and model keep changes inside the given subrange. The model carries the typed panic and trace while suppressing all later comparisons and writes.

My `go run ./difftest/slices/local-sorts` replay passed 12,963 boundary rows across the three pairs and 20,000 sampled rows per pair for each seed 17, 2718 and 65537 against the pinned private Go overlay and compiled Lean. The independent `go run ./difftest/slices/local-sorts/review` passed six literal fixtures, two per pair, and checked insertion rollback, sift comparison order and heap drain exchange mutations. [Independent control](../../../difftest/slices/local-sorts/review/main.go)

The native observations cover valid int64 subranges and representable heap indices, plus comparator scripts that may return inconsistent results or panic. Callback mutation, invalid bounds, signed index overflow, concurrent access, and other generic instantiations are outside this replay.
