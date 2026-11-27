# Independent review of `time.absYday.janFeb`

**Verdict: FAITHFUL**

The source switches from zero to one at March-based day 306. The model’s theorem and the controls check December 31 and absolute January 1.

Within the author scope, absolute seconds, absolute days and century words are uint64, within-century years are 0–99, and March-based year days are 0–365. The scoped Lean audit found no axiom dependencies in the six registered theorems.

Pinned Go 1.27.0 and Lean passed 5,124 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed. Ignoring weekday wrap disagrees on 3, reporting standard month numbers for January and February disagrees on 4, and applying only the four-year leap rule disagrees on 2.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-fields/review
```

Evidence SHA-256: `0b6a4d00e7242b550a94afdbdf540b12ace581efa5ce06e9c3ec86d5cdb193de`
