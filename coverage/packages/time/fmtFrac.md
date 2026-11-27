# Independent review of `time.fmtFrac`

**Verdict: FAITHFUL**

Evidence SHA-256: `1fb178de0a3398e010191c37a188ad4e4a377e70774466a831847973e68c5edc`

I checked evidence commit `9ee1961` against Go 1.27.0 `src/time/time.go`. Pinned Go processes digits from least to most significant, skipping trailing zeros before any output. It writes the decimal point last and only after a digit. The control distinguishes a point-write panic after completed digits, an exact suffix index and quotient, an all-zero fractional result with no writes, and zero or negative precision with the original value returned.

Pinned Go and Lean agreed on 2,592 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/format-helpers`. Seven independent cases passed with `go run ./difftest/time/format-helpers/review` and rejected nine incorrect observations. The scoped kernel audit passed. Scope covers uint64 input, legal visible length and backing capacity, and the bounds failure class. Positive precision above 32 is modeled and outside the replayed domain. [Independent control](../../../difftest/time/format-helpers/review/main.go)
