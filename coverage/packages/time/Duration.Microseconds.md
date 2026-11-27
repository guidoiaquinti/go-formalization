# Independent review of `time.Duration.Microseconds`

**Verdict: FAITHFUL**

Evidence SHA-256: `4b1ed03a53c4acc7c6bcceb86afd79b6434e9af046a2096f89ca123a2ecc25e4`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. The source divides by 1,000 with signed truncation toward zero. The registered theorem establishes the negative boundary at -999, and the independent control distinguishes it from floor division.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
