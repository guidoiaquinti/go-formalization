# Independent review of `time.fmtInt`

**Verdict: FAITHFUL**

Evidence SHA-256: `fbc0d45d664b8951acf77535b15afbba212e5a20b6d3cedf2f13b990c2492912`

I checked evidence commit `9ee1961` against Go 1.27.0 `src/time/time.go`. Pinned Go writes from the end of the visible slice toward the front and returns the start index. Zero emits a single digit. The control confirms that a buffer with only two visible cells retains the last two digits of 567 on a bounds panic and preserves the spare backing byte.

Pinned Go and Lean agreed on 2,592 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/format-helpers`. Seven independent cases passed with `go run ./difftest/time/format-helpers/review` and rejected nine incorrect observations. The scoped kernel audit passed. Scope covers uint64 input, legal visible length and backing capacity, and the bounds failure class. Positive precision above 32 is modeled and outside the replayed domain. [Independent control](../../../difftest/time/format-helpers/review/main.go)
