# Independent review of `time.absYday.split`

**Verdict: FAITHFUL**

The source packs a linear day calculation into a uint32 word then splits high and low bits. The month is absolute, so January and February are 13 and 14. Controls include March 31 and April 1, December to January, and both final February days.

Within the author scope, absolute seconds, absolute days and century words are uint64, within-century years are 0–99, and March-based year days are 0–365. The scoped Lean audit found no axiom dependencies in the six registered theorems.

Pinned Go 1.27.0 and Lean passed 5,124 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed. Ignoring weekday wrap disagrees on 3, reporting standard month numbers for January and February disagrees on 4, and applying only the four-year leap rule disagrees on 2.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-fields/review
```

Evidence SHA-256: `d9d00a4bf0129652100412a37acd26c387b86cb72e1c8fbfc6bf34961efb6714`
