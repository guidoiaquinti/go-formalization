# Independent review of `strconv.FormatUint`

**Verdict: FAITHFUL**

Evidence SHA-256: `5ff8a68aef2dcd65e52a891dd51dcff44d15ee0cfc519cf8fdae6f07999a4bc1`

I reviewed evidence commit `148f1ba` against Go 1.27.0 `src/strconv/number.go` and its delegation to `internal/strconv`. The public wrapper forwards its arguments and observations directly to the internal implementation. The internal source uses specialized decimal and power-of-two paths. The byte-level model yields equivalent unsigned digits using quotient and remainder. The fixed base-36 example rejects uppercase, and native comparisons check the exact invalid-base panic text. Lean represents that panic as a flag.

Pinned native Go, Lean and the independent radix and buffer oracles agreed on 4,382 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64. I replayed `go run ./difftest/strconv/integer-format` and five independent rows with `go run ./difftest/strconv/integer-format/review`. The group rejects twelve deliberately wrong observations across the two packages. Scope covers ordinary slice geometry and observable backing writes. Allocator-selected growth capacity and native 32-bit replay are outside this review. [Independent control](../../../difftest/strconv/integer-format/review/main.go)
