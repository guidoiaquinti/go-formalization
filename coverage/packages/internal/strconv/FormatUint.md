# Independent review of `internal/strconv.FormatUint`

**Verdict: FAITHFUL**

Evidence SHA-256: `196afd05d6a727e532d8078b9babe88cdf0d121c04a0c426ac5c13ee6a1d63ee`

I reviewed evidence commit `148f1ba` against Go 1.27.0 `src/internal/strconv/itoa.go`. The internal source uses specialized decimal and power-of-two paths. The byte-level model yields equivalent unsigned digits using quotient and remainder. The fixed base-36 example rejects uppercase, and native comparisons check the exact invalid-base panic text. Lean represents that panic as a flag.

Pinned native Go, Lean and the independent radix and buffer oracles agreed on 4,382 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64. I replayed `go run ./difftest/strconv/integer-format` and five independent rows with `go run ./difftest/strconv/integer-format/review`. The group rejects twelve deliberately wrong observations across the two packages. Scope covers ordinary slice geometry and observable backing writes. Allocator-selected growth capacity and native 32-bit replay are outside this review. [Independent control](../../../../difftest/strconv/integer-format/review/main.go)
