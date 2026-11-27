# Independent review of `time.absSeconds.clock`

**Verdict: FAITHFUL**

The source reduces by 86,400 before decomposing into hours, minutes and seconds. The final second of a day and first second of the next day distinguish omitted day reduction. Seconds at minute and hour boundaries are also checked.

Within the author scope, absolute seconds, absolute days and century words are uint64, within-century years are 0–99, and March-based year days are 0–365. The scoped Lean audit found no axiom dependencies in the six registered theorems.

Pinned Go 1.27.0 and Lean passed 5,124 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed. Ignoring weekday wrap disagrees on 3, reporting standard month numbers for January and February disagrees on 4, and applying only the four-year leap rule disagrees on 2.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-fields/review
```

Evidence SHA-256: `57cb1a1b9b8bcad5c6b2a8c67bbabd3c12591c04c302dd4979f976fce675797e`
