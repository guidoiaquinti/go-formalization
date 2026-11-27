# Independent review of `time.absDays.yearYday`

**Verdict: FAITHFUL**

The source composes the absolute split with civil year and ordinal calculations. The theorem covers January rollover. Controls verify February 29 in 2000 and 2400 and the distinct March 1 ordinals in 1900 and 2000.

The author scope covers the full unsigned day and century words, with regular calendar subfields for the standalone conversions. Native civil-year conversion is parameterized by int width. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 6,954 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Thirteen independent literal controls passed. Omitting day-word wrap or applying the leap correction before January rollover would disagree with the corresponding controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/calendar-conversion/review
```

Evidence SHA-256: `c9928c0331d6dc5bc20b85e0fe1f4ccf1ad8f6b5ba7f6f4862b52fab49cb1b54`
