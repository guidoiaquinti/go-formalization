# Independent review of `sort.medianAdjacent`

**Verdict: FAITHFUL**

Evidence SHA-256: `e4be03f8bc0ef0f7877e39edad143f683e975451d53ca5c25eac264e4e7dc0ef`

I read the pinned Go 1.27.0 declaration and checked the model, theorem claims and comparison. The network receives the predecessor, center and successor indices, and updates the counter through the same three comparisons. The Go package overlay invokes the exact private symbol, records the ordered Less arguments, and rejects any call to Len or Swap. The comparison covers both callback implementations with scripted replies, including scripts that need not describe a consistent ordering.

I ran `go run ./difftest/sort/pivot-helpers/review` on darwin/arm64 with Go 1.27.0. Eleven independent literal traces matched both native implementations and compiled Lean. Leaving the swap counter unchanged failed three fixtures, one for each helper that takes a counter. Reusing the median-only path at span 50 failed two fixtures. Dropping the decreasing hint failed one. The author comparison independently enumerates 4,096 complete ninther scripts and samples each helper kind with 60,000 rows. [Reviewer control](../../../difftest/sort/pivot-helpers/review/main.go)

The scope is sequential access with valid index ranges, a non-nil comparator that returns normally, and a valid counter pointer where applicable. Counter wrap is modeled on signed 64-bit Go.
