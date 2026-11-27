# Independent review of `slices.swapRangeOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `c3bb1da98f0cc952ed5c679e69e73afb899a790f0e50ed6e99e1fd3ee81e6f3d`

I checked evidence commit `70bc0e3` against the pinned Go 1.27.0 private implementation. Each loop iteration reads the two current cells and swaps them, so the model correctly composes swaps in ascending order even with overlap. Both accesses are evaluated before the stores. If an access faults, all prior swaps survive and the failing pair leaves the backing as it was after those prior swaps. Valid indices start in the interval from zero through length minus one and advance to the first out-of-bounds index before signed addition can overflow. A valid native slice can therefore run at most its visible length in successful swaps. Nonpositive counts skip the loop entirely, and the comparator variant never calls its comparator.

I ran `go run ./difftest/slices/swap-range` and independently reproduced 4,855 boundary cases and 20,000 per seed 17, 2718, and 65537 against the same-package pinned native Go overlay and the compiled Lean model. I also ran `go run ./difftest/slices/swap-range/review`, adding eight literal fixtures with both engines. On backing `[4,7,1,8,91]`, overlapping forward swaps produce `[7,1,8,4,91]`; the checked descending-order mutation produces `[8,4,7,1,91]`. With a visible length of four and offsets zero and two, the third pair faults after producing `[1,8,4,7,91]`. The preflight mutation rejects before changing anything, and the capacity-for-length mutation changes the spare cell to produce `[1,8,91,7,4]`. Nil, empty, identical locations, signed extrema, and panic-on-call comparators were also checked. [Independent control](../../../difftest/slices/swap-range/review/main.go)

The scope covers sequential access with legal slice geometry, signed native offsets and count, and index-bounds panic classification. The native overlay instantiates the generic declarations with `int` on darwin/arm64. Exact runtime error text and concurrent access are outside scope.
