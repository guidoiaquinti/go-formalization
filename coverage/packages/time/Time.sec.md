# Independent review of `time.Time.sec`

**Verdict: FAITHFUL**

The monotonic branch strips the flag, extracts the 33-bit seconds field and adds the internal epoch offset for 1885. The other branch returns signed ext, as proved for every unmarked word. Controls include both packed seconds endpoints and the minimum signed ext.

The review covers raw wall and ext words in the pinned Go 1.27.0 declarations, including strictly synthetic words and fixed native locations. It exercises 18,769 fixed rows, 60,000 sampled rows at seeds 17, 2718 and 65537, and 9 independent literal controls on darwin/arm64. Private methods run through a test-only bridge with nonnil receivers. The scoped Lean axiom audit found only permitted dependencies.

The shared reviewer control distinguishes four mutations. Decoding zero by instant coordinates disagrees on 1 row. Unsigned monotonic comparisons disagree on 1 row. Always comparing wall coordinates for marked pairs disagrees on 2 rows, and choosing ext when either operand is marked disagrees on 2 rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/time-compare/review
```

Evidence SHA-256: `e5416270087ffe53682328e907551039013b94e7bf66dc2827f78b90fa6214d6`
