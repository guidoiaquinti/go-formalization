# Independent review of `time.Time.Compare`

**Verdict: FAITHFUL**

The source picks the pair of signed monotonic readings only when both operands are marked, otherwise comparing wall seconds and then nanoseconds. It returns one of -1, 0, 1. The theorem checks the two-marked branch and controls exercise all three signs.

The review covers raw wall and ext words in the pinned Go 1.27.0 declarations, including strictly synthetic words and fixed native locations. It exercises 18,769 fixed rows, 60,000 sampled rows at seeds 17, 2718 and 65537, and 9 independent literal controls on darwin/arm64. Private methods run through a test-only bridge with nonnil receivers. The scoped Lean axiom audit found only permitted dependencies.

The shared reviewer control distinguishes four mutations. Decoding zero by instant coordinates disagrees on 1 row. Unsigned monotonic comparisons disagree on 1 row. Always comparing wall coordinates for marked pairs disagrees on 2 rows, and choosing ext when either operand is marked disagrees on 2 rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/time-compare/review
```

Evidence SHA-256: `2ae1c3dccb46d546d10d411eab05962d372e41ec32d73d1e520bad75cdfb6125`
