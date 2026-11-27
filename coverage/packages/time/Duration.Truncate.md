# Independent review of `time.Duration.Truncate`

**Verdict: FAITHFUL**

Evidence SHA-256: `a0b321befbb65aa07fe7919c3fd8edfaffc2f3feceadf42bdfb49db995245ad9`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. The positive-step source subtracts the signed remainder. The theorem covers the nonpositive guard. The independent negative example distinguishes truncation toward zero from floor rounding.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
