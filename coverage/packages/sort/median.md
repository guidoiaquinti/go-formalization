# Independent review of `sort.median`

**Verdict: FAITHFUL**

Evidence SHA-256: `d9df39748b4a65d555d35499f463625789a5a2ee7d92e2038990c83ce7b4b897`

I read the pinned Go 1.27.0 declaration and checked the model, theorem claims and comparison. The three comparisons run in network order against the currently selected indices. The result is the middle index after the network, and the count changes by the number of true replies. The Go package overlay invokes the exact private symbol, records the ordered Less arguments, and rejects any call to Len or Swap. The comparison covers both callback implementations with scripted replies, including scripts that need not describe a consistent ordering.

I ran `go run ./difftest/sort/pivot-helpers/review` on darwin/arm64 with Go 1.27.0. Eleven independent literal traces matched both native implementations and compiled Lean. Leaving the swap counter unchanged failed three fixtures, one for each helper that takes a counter. Reusing the median-only path at span 50 failed two fixtures. Dropping the decreasing hint failed one. The author comparison independently enumerates 4,096 complete ninther scripts and samples each helper kind with 60,000 rows. [Reviewer control](../../../difftest/sort/pivot-helpers/review/main.go)

The scope is sequential access with valid index ranges, a non-nil comparator that returns normally, and a valid counter pointer where applicable. Counter wrap is modeled on signed 64-bit Go.
