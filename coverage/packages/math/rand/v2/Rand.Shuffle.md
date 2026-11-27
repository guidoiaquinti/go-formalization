# Independent review of `math/rand/v2.Rand.Shuffle`

**Verdict: FAITHFUL**

Evidence SHA-256: `546c25534dc78cc38a9329ba607a7b7651673fccb29c72a7c65ffcf4d99a48e6`

Visits indices from n−1 down to one, retaining each prior invocation and consuming rejected words before invoking the next swap. The callback panic retains its exact payload and attempted arguments, stops further draws, and occurs after any source rejection retries. A nil callback is unused for zero and one elements and can only panic after drawing for a larger input. All 65,760 author cases passed, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Shuffle.lean) · [Comparison](../../../../difftest/math/rand/v2/shuffle/check.go)

Eight independent native and Lean cases passed, with six altered observations rejected. The fixed bound-three draw rejects an initial zero and accepts a low product exactly at the rejection threshold. The scope uses feasible native lengths and deterministic sources and callbacks. [Reviewer control](../../../../difftest/math/rand/v2/shuffle/review/main.go)
