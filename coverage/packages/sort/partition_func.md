# Independent review of `sort.partition_func`

**Verdict: FAITHFUL**

Evidence SHA-256: `3696aa9e502516091fa43aefbbbbf585983a820aeefc19ccf688b1395695df8c`

I traced the exact pinned Go 1.27.0 declaration alongside the Lean model and compared the test hypotheses with the source. The pinned declaration installs the selected pivot at the start, scans inward under a strict less-than comparison, and moves the pivot into its final slot. A crossing before any interior swap produces the already-partitioned flag. The singleton case still issues two self-swaps.

I ran `go run ./difftest/sort/partition/review` on darwin/arm64. Ten independent literal traces matched native Go and compiled Lean. The witnesses reject an always-true partition flag, an omitted singleton self-swap, and an off-by-one equal boundary across both callback variants. The author compared 10,435 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 for this declaration. [Reviewer control](../../../difftest/sort/partition/review/main.go)

The tested domain is a nonempty in-bounds range and a valid pivot. Less and Swap operate on signed integer array elements and return normally. A minimum-valued pivot is required for partitionEqual.
