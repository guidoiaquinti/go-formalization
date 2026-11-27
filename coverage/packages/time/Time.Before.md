# Independent review of `time.Time.Before`

**Verdict: FAITHFUL**

The strict ascending comparison uses signed ext only when both operands are marked. Otherwise it uses wall seconds then nanoseconds. The theorem checks the two-marked branch and controls include mixed flags and nanosecond ties.

The review covers raw wall and ext words in the pinned Go 1.27.0 declarations, including strictly synthetic words and fixed native locations. It exercises 18,769 fixed rows, 60,000 sampled rows at seeds 17, 2718 and 65537, and 9 independent literal controls on darwin/arm64. Private methods run through a test-only bridge with nonnil receivers. The scoped Lean axiom audit found only permitted dependencies.

The shared reviewer control distinguishes four mutations. Decoding zero by instant coordinates disagrees on 1 row. Unsigned monotonic comparisons disagree on 1 row. Always comparing wall coordinates for marked pairs disagrees on 2 rows, and choosing ext when either operand is marked disagrees on 2 rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/time-compare/review
```

Evidence SHA-256: `504aa798b3f90890ca15df5b46e0d87a68ee31576f67d10d391421a415a54541`
