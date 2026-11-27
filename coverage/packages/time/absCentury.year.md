# Independent review of `time.absCentury.year`

**Verdict: FAITHFUL**

The source converts the century offset to native signed int and adds the within-century year and January/February rollover. The model applies the equivalent width reduction after additions. Native controls include maximum uint64 century and a January rollover from year 99.

The author scope covers the full unsigned day and century words, with regular calendar subfields for the standalone conversions. Native civil-year conversion is parameterized by int width. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 6,954 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Thirteen independent literal controls passed. Omitting day-word wrap or applying the leap correction before January rollover would disagree with the corresponding controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/calendar-conversion/review
```

Evidence SHA-256: `ec88b456679f9107e1d090d9731d53758d99a77aafca7ccbd9bcf0a7993b2633`
