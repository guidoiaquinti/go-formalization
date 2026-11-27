# Independent review of `time.absDays.split`

**Verdict: FAITHFUL**

The source uses uint64 scaling and then splits the century remainder through the uint32 multiply-high/low method. The model preserves the scaling wrap and equivalent quotient and remainder. Controls pair day zero with day 2^62 and check the 146,097-day cycle.

The author scope covers the full unsigned day and century words, with regular calendar subfields for the standalone conversions. Native civil-year conversion is parameterized by int width. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 6,954 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Thirteen independent literal controls passed. Omitting day-word wrap or applying the leap correction before January rollover would disagree with the corresponding controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/calendar-conversion/review
```

Evidence SHA-256: `e69aba67e510dd8fbb176438c184da3adb5303e918bc0efa91a9798ae8566000`
