# Independent review of `slices.insertionSortOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `2f873e82588c77a09eda857685ec4ffe2dbdb25c5967eb6966504ec337a6dd2f`

I checked evidence commit `d508653` against the private Go 1.27.0 helpers and their Lean counterparts. Insertion walks outward indices from a+1 and swaps adjacent strict inversions back toward a. The independent descending-subrange fixture preserves exactly the first completed swap when the second comparison panics, and a singleton skips the comparator. The source and model keep changes inside the given subrange. The model carries the typed panic and trace while suppressing all later comparisons and writes.

My `go run ./difftest/slices/local-sorts` replay passed 12,963 boundary rows across the three pairs and 20,000 sampled rows per pair for each seed 17, 2718 and 65537 against the pinned private Go overlay and compiled Lean. The independent `go run ./difftest/slices/local-sorts/review` passed six literal fixtures, two per pair, and checked insertion rollback, sift comparison order and heap drain exchange mutations. [Independent control](../../../difftest/slices/local-sorts/review/main.go)

The native observations cover valid int64 subranges and representable heap indices, plus comparator scripts that may return inconsistent results or panic. Callback mutation, invalid bounds, signed index overflow, concurrent access, and other generic instantiations are outside this replay.
