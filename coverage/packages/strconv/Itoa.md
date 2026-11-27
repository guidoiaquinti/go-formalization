# Independent review of `strconv.Itoa`

**Verdict: FAITHFUL**

Evidence SHA-256: `83388525c1f596bb7370769a9a48f5d7f6ab26c3a2e58a40a0872449769fb1f8`

I reviewed evidence commit `148f1ba` against Go 1.27.0 `src/strconv/number.go` and its delegation to `internal/strconv`. The public wrapper forwards its arguments and observations directly to the internal implementation. The conversion extends a signed native-width integer to int64 and formats it in decimal. The review uses 0xffffffff to distinguish the 32- and 64-bit interpretations when replayed at those widths. Native replay here was at 64 bits.

Pinned native Go, Lean and the independent radix and buffer oracles agreed on 4,382 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64. I replayed `go run ./difftest/strconv/integer-format` and five independent rows with `go run ./difftest/strconv/integer-format/review`. The group rejects twelve deliberately wrong observations across the two packages. Scope covers ordinary slice geometry and observable backing writes. Allocator-selected growth capacity and native 32-bit replay are outside this review. [Independent control](../../../difftest/strconv/integer-format/review/main.go)
