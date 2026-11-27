# Independent review of `time.Duration.Seconds`

**Verdict: FAITHFUL**

Evidence SHA-256: `c20ebffd6bdc2f45cac3a23d1142191054aabd3b6d9933ae34775d16d6100755`

I checked evidence commit `c98bc74` against pinned Go 1.27.0 `src/time/time.go`. Go splits the count into signed whole seconds and remainder before floating conversion and addition. The quotient, remainder and divisor are exactly representable in binary64 over the int64 domain. The theorem states the two-stage expression and the rational differential oracle checks exact rounding. At -814638300595526980 nanoseconds, converting the full input before dividing changes the binary64 result.

The direct pinned Go and Lean comparison and independent oracles passed 712 deterministic cases and 20,000 cases per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/time/duration-numeric`. Nine independently selected fixtures passed with `go run ./difftest/time/duration-numeric/review`, and twelve wrong observations across the group were rejected. The modeled input scope is signed int64 durations and divisors. [Independent control](../../../difftest/time/duration-numeric/review/main.go)
