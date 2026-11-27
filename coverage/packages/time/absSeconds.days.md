# Independent review of `time.absSeconds.days`

**Verdict: FAITHFUL**

The source divides the unsigned absolute seconds by 86,400. The kernel theorem and controls check the day transition at 86,399 and 86,400 and the top uint64 value.

Within the author scope, absolute seconds, absolute days and century words are uint64, within-century years are 0–99, and March-based year days are 0–365. The scoped Lean audit found no axiom dependencies in the six registered theorems.

Pinned Go 1.27.0 and Lean passed 5,124 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed. Ignoring weekday wrap disagrees on 3, reporting standard month numbers for January and February disagrees on 4, and applying only the four-year leap rule disagrees on 2.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-fields/review
```

Evidence SHA-256: `2ea0bd14a5793a5374cea4f0c52ccfa504bdf8d878e56e34669801bcfb5ec979`
