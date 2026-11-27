# Independent review of `sort.partialInsertionSort`

**Verdict: FAITHFUL**

Evidence SHA-256: `050bb9bd034aa12deb190bfd7d8aebdb238ebcfcacba0331c2bea0d7599a36fe`

I checked the pinned Go 1.27.0 loop, the Lean state transition and the exact private native entry point. An empty selected range returns false. A short unsorted range returns false without swapping. For longer ranges, the loop can repair up to five inversions and follows the source bounds for each shift.

I ran `go run ./difftest/sort/partial-insertion-sort/review` on darwin/arm64. Eight independent exact traces matched native Go and compiled Lean. The checks distinguish a generic empty-range sortedness result and a left shift limited to the selected range. In a length-54 input with prefix `[100,101,0,1,-1]` and range `[2,52)`, both helpers return true and change the prefix to `[-1,100,101,0,1]`. The trace includes swaps at indices 2 and 1, then 1 and 0. The author replay covered 6,813 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 per declaration. [Reviewer control](../../../difftest/sort/partial-insertion-sort/review/main.go)

The tested domain covers in-bounds signed integer array ranges and ordinary indexed Less and Swap callbacks. The offset fixture describes the private helper. Public caller reachability remains unverified.
