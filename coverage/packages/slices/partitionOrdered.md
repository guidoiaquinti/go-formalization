# Independent review of `slices.partitionOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `46f8ef0006e2e19e59ee4ba2b65d778a9431357098984c7a7dd1786859d893ca`

I checked evidence commit `3ea0b44` against the private Go 1.27.0 source. Both helpers move the selected pivot to the front before comparing and scan inclusive remaining endpoints in left-then-right order. Ordinary partition places the pivot at the right stop after the scans cross and reports already partitioned precisely when it did no cross-swap. Equal partition keeps the pivot at the front, uses it as the first comparison operand on both scans, and returns the first index after the accepted group. That group is equal to the pivot when the pivot is minimal. A callback panic preserves the initial pivot move and completed cross-swaps and skips the pending writes.

I ran `go run ./difftest/slices/partition` and reproduced 6,482 boundaries per pair and 20,000 cases per pair for each seed 17, 2718 and 65537 on pinned native Go and compiled Lean. Six independent literal fixtures also passed with `go run ./difftest/slices/partition/review`. They detect counting pivot-only movement as a cross-swap, rolling back a completed cross-swap on a later panic, and reversing the first comparator operands in equal partition. The controls also cover a pivot above the minimum, duplicate pivot values, unchanged guards and no callback on a singleton. [Independent control](../../../difftest/slices/partition/review/main.go)

The scope covers valid int64 ranges with the pivot inside a nonempty interval and scripted, non-mutating callbacks. Invalid bounds, arbitrary callback mutation and concurrent access are outside the claim.
