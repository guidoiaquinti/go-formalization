# Independent review of `time.Duration.Abs`

**Verdict: FAITHFUL**

Evidence SHA-256: `39a1fd5d0e9e81363f4ae3c655b8058839e78c7661644d26fbdc102d79fb74ea`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. The source special-cases the minimum signed value, negates other negative values and leaves nonnegative values unchanged. Lean models unsigned magnitude clamped at the maximum. The minimum theorem and independent native control distinguish saturation from signed negation that wraps.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
