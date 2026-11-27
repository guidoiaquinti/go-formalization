# Independent review of `internal/strconv.FormatInt`

**Verdict: FAITHFUL**

Evidence SHA-256: `58f5a2c67aa2e52a038dbb2c2fc6bdd2015b7c3f162c97c4488d78f23ee11637`

I reviewed evidence commit `148f1ba` against Go 1.27.0 `src/internal/strconv/itoa.go`. The positive magnitude fits in an unsigned word even for the minimum signed input. The byte-level model preserves its leading minus sign and native tests compare the exact pinned invalid-base panic text. Lean represents that panic as a flag.

Pinned native Go, Lean and the independent radix and buffer oracles agreed on 4,382 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64. I replayed `go run ./difftest/strconv/integer-format` and five independent rows with `go run ./difftest/strconv/integer-format/review`. The group rejects twelve deliberately wrong observations across the two packages. Scope covers ordinary slice geometry and observable backing writes. Allocator-selected growth capacity and native 32-bit replay are outside this review. [Independent control](../../../../difftest/strconv/integer-format/review/main.go)
