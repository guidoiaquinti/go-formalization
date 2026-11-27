# Independent review of `time.absYday.yday`

**Verdict: FAITHFUL**

The source converts to a one-based ordinal and applies the leap bit only for March through December. Fixed controls distinguish March 1 in 1900, 2000 and 2100 and confirm that February 29 keeps ordinal 60.

The author scope covers the full unsigned day and century words, with regular calendar subfields for the standalone conversions. Native civil-year conversion is parameterized by int width. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 6,954 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Thirteen independent literal controls passed. Omitting day-word wrap or applying the leap correction before January rollover would disagree with the corresponding controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/calendar-conversion/review
```

Evidence SHA-256: `f717eef9721067989d83cb69d106e46128e34ba3dbcaa0fd8c9850a7e396453d`
