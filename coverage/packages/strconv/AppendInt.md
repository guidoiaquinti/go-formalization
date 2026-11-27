# Independent review of `strconv.AppendInt`

**Verdict: FAITHFUL**

Evidence SHA-256: `5140b9255dd18d77648b1e7260aedd5d15d738475936c7c730c2b1cf3f43ac12`

I reviewed evidence commit `148f1ba` against Go 1.27.0 `src/strconv/number.go` and its delegation to `internal/strconv`. The public wrapper forwards its arguments and observations directly to the internal implementation. The negative sign is appended before radix validation and before the digit append. A destination with one spare byte makes the sign observable through the old backing both when invalid radix panics and when the digits force growth.

Pinned native Go, Lean and the independent radix and buffer oracles agreed on 4,382 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64. I replayed `go run ./difftest/strconv/integer-format` and five independent rows with `go run ./difftest/strconv/integer-format/review`. The group rejects twelve deliberately wrong observations across the two packages. Scope covers ordinary slice geometry and observable backing writes. Allocator-selected growth capacity and native 32-bit replay are outside this review. [Independent control](../../../difftest/strconv/integer-format/review/main.go)
