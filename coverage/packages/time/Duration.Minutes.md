# Independent review of `time.Duration.Minutes`

**Verdict: FAITHFUL**

Evidence SHA-256: `9752c0e88952d2db17e93ed921d13d33dc89e0a97b6a42ff7576a061909a598c`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. Go splits the count into signed whole minutes and remainder before floating conversion and addition. The quotient, remainder and divisor are exactly representable in binary64 over the int64 domain. The theorem states the two-stage expression and the rational differential oracle checks exact rounding. At -814638300595526980 nanoseconds, converting the full input before dividing changes the binary64 result.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
