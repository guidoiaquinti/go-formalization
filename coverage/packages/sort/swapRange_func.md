# Independent review of `sort.swapRange_func`

**Verdict: FAITHFUL**

Evidence SHA-256: `db6f8dfab24642f34fdb3f2c8eb8e2b1e331c34821e3f31b4efd532e9d604f31`

I checked the Go 1.27.0 private implementation against the Lean trace model and the independent permutation oracle. The helper walks ascending offsets and invokes every swap including self-swaps. With overlapping windows it composes sequential transpositions. Reversing the order fails on the overlapping fixture. The native overlay calls the exact declaration, records every Swap, and panics if Len or Less is called.

The author replay passed 3,185 exhaustive geometries through width 12 and 60,000 seeded cases (17, 2718 and 65537) for this declaration. I ran `go run ./difftest/sort/range-moves/review` on darwin/arm64. Sixteen independent literal traces matched native Go and Lean. All six variant-specific mutants were rejected. [Reviewer control](../../../difftest/sort/range-moves/review/main.go)

Scope is valid index geometry with nonpanicking callbacks. Rotation requires two nonempty blocks or the all-empty case. Exactly one empty block is outside the pinned caller domain.
