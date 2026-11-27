# Independent review of `slices.choosePivotCmpFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `1e7c8d88fd57e2276a983922a9d86a321a4355a4cbe60f6c628f3793d66c1365`

I checked evidence commit `29d1fe4` against the private Go 1.27.0 helper. The model computes each quarter index from the same truncated quarter length, skips comparisons below length eight, and uses three adjacent medians before the final median at length fifty or more. The shared counter records all exchanges. Zero exchanges produce the increasing hint, twelve produce decreasing, and other counts produce unknown. The comparator case hands each nested median the accumulated invocation count and preserves exactly the completed callback trace when a call panics. The backing remains intact.

My `go run ./difftest/slices/choose-pivot` replay passed 22,464 boundaries and 20,000 cases per seed 17, 2718 and 65537 against pinned native Go and compiled Lean on darwin/arm64. Seven independent literal fixtures also passed both with `go run ./difftest/slices/choose-pivot/review`. For range [3,10), an ordinary midpoint returns 6 where the pinned quarter-based position is 5. Delaying the ninther past length fifty starts with pair (25,13) where the real first pair is (13,12). A descending length-49 witness exchanges three pairs and reports unknown, whereas length 50 reports decreasing after twelve exchanges. The control also checks no callback below eight, the three-call threshold at eight, all twelve successful calls at fifty, and a typed panic in the second neighborhood. [Independent control](../../../difftest/slices/choose-pivot/review/main.go)

The native observations cover valid ranges on int64 backing and scripted non-mutating callbacks. Invalid bounds, arbitrary callback mutation and concurrent access are outside the registered scope.
