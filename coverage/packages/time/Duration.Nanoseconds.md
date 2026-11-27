# Independent review of `time.Duration.Nanoseconds`

**Verdict: FAITHFUL**

Evidence SHA-256: `c910ced660b6831361db62cc92214d5658cfc9126c5bf41b96e900c952ae0c48`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. The source returns the signed int64 representation directly. The identity theorem covers the modeled domain. The minimum duration rejects exposing the raw word as an unsigned result.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
