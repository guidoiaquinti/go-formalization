# Independent review of `slices.pdqsortOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `7c591e9ffc950028cf4297712a9fe238200c7f2c12b2f3a760abb63967f3bd2f`

I checked evidence commit `c759924` against the private Go 1.27.0 source. The model checks the insertion limit before the heap fallback, decrements the limit at an unbalanced iteration, and uses the pinned pivot, optional reversal and partial repair, and predecessor-conditioned equal partition. It resets both optimization flags in each recursive call and carries the correct balance and partition flags into the larger tail range. A valid partition removes at least one element along any recursive path, which justifies the model's fuel. The comparator variant carries one global callback trace into subcalls and preserves committed mutations if a comparison panics.

I ran `go run ./difftest/slices/pdqsort` and reproduced 38,016 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go and compiled Lean. The native check also asserts full multiset preservation, exact suffix, and sorted output and prefix preservation where the predecessor invariant applies. Five independent literal fixtures passed with `go run ./difftest/slices/pdqsort/review`. They distinguish insertion before fallback at length twelve, zero-limit heap dispatch at length thirteen, and reversal before the first post-ninther comparison at length fifty. The last fixture panics at callback 12 and observes the completed reversal plus the pre-panic trace. [Independent control](../../../difftest/slices/pdqsort/review/main.go)

The native replay is darwin/arm64 with int64 backing, valid ranges, representable index arithmetic and fixed non-mutating comparator scripts. Invalid bounds, arbitrary callback mutation and concurrent access are outside scope.
