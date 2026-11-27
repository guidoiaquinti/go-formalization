# Independent review of `time.Duration.Round`

**Verdict: FAITHFUL**

Evidence SHA-256: `8a410b461d89ff4b40a7ecbf4f3f4185e9fd8d2f62f545b4fc552ee7ad44d373`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. The source compares the absolute remainder to half the divisor without overflowing and detects overflow of the signed addition or subtraction. Lean uses unbounded arithmetic and clamps to the signed endpoints. The equivalence is checked by differential comparisons. The registered theorems cover the nonpositive guard and concrete tie and saturation edges. Independent cases distinguish both tie signs, overflow wrapping, and a below-half remainder with an odd divisor.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
