# Independent review of `time.Time.unixSec`

**Verdict: FAITHFUL**

The source adds the negative internal-to-Unix offset to signed wall seconds in int64. The model preserves the width and the controls include minimum-ext wraparound and the Unix epoch. The theorem establishes the epoch result.

The review covers raw wall and ext words in the pinned Go 1.27.0 declarations, including strictly synthetic words and fixed native locations. It exercises 18,769 fixed rows, 60,000 sampled rows at seeds 17, 2718 and 65537, and 9 independent literal controls on darwin/arm64. Private methods run through a test-only bridge with nonnil receivers. The scoped Lean axiom audit found only permitted dependencies.

The shared reviewer control distinguishes four mutations. Decoding zero by instant coordinates disagrees on 1 row. Unsigned monotonic comparisons disagree on 1 row. Always comparing wall coordinates for marked pairs disagrees on 2 rows, and choosing ext when either operand is marked disagrees on 2 rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/time-compare/review
```

Evidence SHA-256: `b0200bc237a3838d874829b7d857147d22e12c4d113818c35ca1ac28a1eac710`
