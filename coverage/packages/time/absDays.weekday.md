# Independent review of `time.absDays.weekday`

**Verdict: FAITHFUL**

The source adds the Wednesday offset in uint64 before reducing modulo seven. The exact three input words that wrap are all in the reviewer controls. Reducing before adding changes those answers.

Within the author scope, absolute seconds, absolute days and century words are uint64, within-century years are 0–99, and March-based year days are 0–365. The scoped Lean audit found no axiom dependencies in the six registered theorems.

Pinned Go 1.27.0 and Lean passed 5,124 fixed rows and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed. Ignoring weekday wrap disagrees on 3, reporting standard month numbers for January and February disagrees on 4, and applying only the four-year leap rule disagrees on 2.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-fields/review
```

Evidence SHA-256: `57737579b2efab042e9b36052afac15a493c083fed86d3e6cb441d3929d41b9c`
