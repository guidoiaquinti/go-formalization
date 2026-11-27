# Independent review of `sort.breakPatterns_func`

**Verdict: FAITHFUL**

Evidence SHA-256: `2b04185d1571567b2a79d6bcd4f35f73a6fe4d6b32f14bd845fbcffe13826c59`

I inspected the pinned Go 1.27.0 implementation and checked the Lean model against its arithmetic and callback trace. The function adapter follows the same three-swap trace as the interface variant and performs no swaps below span eight. Starting at 19 with span nine, taking random modulo length changes the recorded callback destinations.

I ran `go run ./difftest/sort/break-patterns/review` on darwin/arm64. Thirteen independent cases matched compiled Lean and the exact pinned Go declarations. Six mutation witnesses failed as expected. The author comparison matched 1,260 boundaries and 60,000 samples at seeds 17, 2718 and 65537 for this symbol. [Reviewer control](../../../difftest/sort/break-patterns/review/main.go)

The modeled ranges are ordered and nonnegative. Callback panics and a nil RNG receiver are outside scope. Native 32-bit execution remains unverified.
