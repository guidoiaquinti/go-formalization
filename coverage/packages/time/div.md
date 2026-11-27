# Independent review of `time.div`

**Verdict: FAITHFUL**

Positive divisors exercise the exact subsecond fast path, the whole-second path and long division. The latter uses the unsigned magnitude assembled from the two pinned words. The negative nonmultiple transformation inverts quotient parity and complements the remainder. Fixed controls include exact divisibility at signed minimum seconds, one nanosecond before the zero instant, and a long-division case outside both fast paths.

The scope follows the pinned packed clock words and preserves opaque location identity. The private division helper is exercised only for positive durations. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 2,916 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent literal controls passed on both implementations, covering strict-half comparisons, positive and negative ties, parity and remainder, and monotonic stripping.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-round/review
```

Evidence SHA-256: `676368d0977ed76a611b84bfd1a2d38ac98fce5d2a2d9a610d066b566a7e4839`
