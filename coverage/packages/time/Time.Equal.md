# Independent review of `time.Time.Equal`

**Verdict: FAITHFUL**

Two marked values can be equal with different wall coordinates if their signed monotonic readings match. A mixed pair instead compares wall coordinates. The theorem checks the two-marked branch and fixed controls also vary receiver locations.

The review covers raw wall and ext words in the pinned Go 1.27.0 declarations, including strictly synthetic words and fixed native locations. It exercises 18,769 fixed rows, 60,000 sampled rows at seeds 17, 2718 and 65537, and 9 independent literal controls on darwin/arm64. Private methods run through a test-only bridge with nonnil receivers. The scoped Lean axiom audit found only permitted dependencies.

The shared reviewer control distinguishes four mutations. Decoding zero by instant coordinates disagrees on 1 row. Unsigned monotonic comparisons disagree on 1 row. Always comparing wall coordinates for marked pairs disagrees on 2 rows, and choosing ext when either operand is marked disagrees on 2 rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/time-compare/review
```

Evidence SHA-256: `68638ed3639e372fab2549809c99853b2dbe2fec2f8615fb507360b22735927b`
