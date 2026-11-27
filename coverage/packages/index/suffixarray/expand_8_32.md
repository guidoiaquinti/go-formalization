# Independent review of `index/suffixarray.expand_8_32`

**Verdict: FAITHFUL**

Evidence SHA-256: `ecf86f41b4a1685665e7b345cefb2a9cbe0cc007feb5ae6ad5c91b44b7ae3ce9`

I checked evidence commit `95f1348` against the Go 1.27.0 SA-IS expansion loops. The byte specialization uses all 256 buckets and signed 32-bit suffix cells. The control uses same-symbol LMS entries and a high raw byte. The backward sweep takes each compacted LMS entry and reserves the highest remaining slot in its symbol bucket. It preserves the relative order of entries within a bucket and clears the other suffix cells.

Pinned native Go and Lean agreed on 8,178 exhaustive ternary-profile rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/index/suffixarray/expand-lms`. Four independent rows passed with `go run ./difftest/index/suffixarray/expand-lms/review` and rejected seven wrong field expectations across the variants. The controls cover equal-symbol placement, zeroed holes, nil or cold or warm frequency state, and preservation of spare backing. The scoped kernel audit passed. Scope assumes a nonempty lawful sorted LMS prefix, disjoint buffers and representable indexes. The ranker in the oracle is fixture generation. Sorting a whole suffix array and malformed intermediate states are outside scope. [Independent control](../../../../difftest/index/suffixarray/expand-lms/review/main.go)
