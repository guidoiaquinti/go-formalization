# Independent review of `time.absDays.date`

**Verdict: FAITHFUL**

The source composes the absolute split, January/February rollover and month conversion. Controls align the start of year 1, the Unix epoch and both leap and nonleap centuries as Gregorian calendar dates.

The author scope covers the full unsigned day and century words, with regular calendar subfields for the standalone conversions. Native civil-year conversion is parameterized by int width. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 6,954 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Thirteen independent literal controls passed. Omitting day-word wrap or applying the leap correction before January rollover would disagree with the corresponding controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/calendar-conversion/review
```

Evidence SHA-256: `7a3af6fef27eae6955dcdaa8d3a5ddc505e82267286ff167aeefdcace8129f14`
