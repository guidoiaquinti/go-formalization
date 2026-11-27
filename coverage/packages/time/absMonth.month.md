# Independent review of `time.absMonth.month`

**Verdict: FAITHFUL**

The source subtracts twelve only for January and February. The theorem and controls distinguish the absolute month numbers 13 and 14 from the civil months 1 and 2.

The author scope covers the full unsigned day and century words, with regular calendar subfields for the standalone conversions. Native civil-year conversion is parameterized by int width. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 6,954 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Thirteen independent literal controls passed. Omitting day-word wrap or applying the leap correction before January rollover would disagree with the corresponding controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/calendar-conversion/review
```

Evidence SHA-256: `b110da58fea6f7ed3bfc8038ebd6a2cc1b0716c9bdb371d84773887c3b043703`
