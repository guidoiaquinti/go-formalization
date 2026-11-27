# Independent review of `math/rand/v2.Rand.Perm`

**Verdict: FAITHFUL**

Evidence SHA-256: `e51021526ab96bc2cf1c1f9f76fc11eaf26504be7a5dcb6ea753cf09d0d7694f`

Constructs a fresh identity list and applies the descending bounded swaps. An empty successful result is nonnil. Source exhaustion prevents a returned permutation, and a negative length panics before drawing. All 65,760 author cases passed, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Shuffle.lean) · [Comparison](../../../../difftest/math/rand/v2/shuffle/check.go)

Eight independent native and Lean cases passed, with six altered observations rejected. The fixed bound-three draw rejects an initial zero and accepts a low product exactly at the rejection threshold. The scope uses feasible native lengths and deterministic sources and callbacks. [Reviewer control](../../../../difftest/math/rand/v2/shuffle/review/main.go)
