# Independent review of `internal/strconv.AppendInt`

**Verdict: FAITHFUL**

Evidence SHA-256: `6dd34a0b822b6c3bacd05f8c5c20a6c64f1998df82de4542fc7ee6930418cd00`

I reviewed evidence commit `148f1ba` against Go 1.27.0 `src/internal/strconv/itoa.go`. The negative sign is appended before radix validation and before the digit append. A destination with one spare byte makes the sign observable through the old backing both when invalid radix panics and when the digits force growth.

Pinned native Go, Lean and the independent radix and buffer oracles agreed on 4,382 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64. I replayed `go run ./difftest/strconv/integer-format` and five independent rows with `go run ./difftest/strconv/integer-format/review`. The group rejects twelve deliberately wrong observations across the two packages. Scope covers ordinary slice geometry and observable backing writes. Allocator-selected growth capacity and native 32-bit replay are outside this review. [Independent control](../../../../difftest/strconv/integer-format/review/main.go)
