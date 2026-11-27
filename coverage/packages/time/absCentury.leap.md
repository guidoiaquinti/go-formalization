# Independent review of `time.absCentury.leap`

**Verdict: FAITHFUL**

For years 0 through 99, ordinary fourth years are leap years. Century year zero is a leap year only in centuries divisible by four. The controls exercise centuries 1 and 4 and a maximum-word century.

Within the author scope, absolute seconds, absolute days and century words are uint64, within-century years are 0–99, and March-based year days are 0–365. The scoped Lean audit found no axiom dependencies in the six registered theorems.

Pinned Go 1.27.0 and Lean passed 5,124 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed. Ignoring weekday wrap disagrees on 3, reporting standard month numbers for January and February disagrees on 4, and applying only the four-year leap rule disagrees on 2.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-fields/review
```

Evidence SHA-256: `60e406e7c5b2a1ee7cf8bf405af27f3f89e93abbf44c478a835af92443b6ee75`
