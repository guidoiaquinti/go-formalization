# Independent review of `sort.breakPatterns`

**Verdict: FAITHFUL**

Evidence SHA-256: `a858494f2a3f45f73c19706ce34c2344727aee62fcc5fa275c92cdd28a248349`

I inspected the pinned Go 1.27.0 implementation and checked the Lean model against its arithmetic and callback trace. Spans below eight have no callbacks. Larger spans seed xorshift with the span and issue three ordered swaps. Starting at 19 with span nine, taking random modulo length selects the wrong second and third destinations.

I ran `go run ./difftest/sort/break-patterns/review` on darwin/arm64. Thirteen independent cases matched compiled Lean and the exact pinned Go declarations. Six mutation witnesses failed as expected. The author comparison matched 1,260 boundaries and 60,000 samples at seeds 17, 2718 and 65537 for this symbol. [Reviewer control](../../../difftest/sort/break-patterns/review/main.go)

The modeled ranges are ordered and nonnegative. Callback panics and a nil RNG receiver are outside scope. Native 32-bit execution remains unverified.
