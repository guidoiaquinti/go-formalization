# Independent review of `net/url.stringContainsCTLByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `8935714bc6cf1353cd49ff9594a1f51cbb8b82d212b59ceeb82e457fc147fe39`

Flags C0 and DEL, while raw 0x80 and 0xFF pass. Independent controls cover NUL, DEL and malformed UTF-8 without classifying high bytes as control bytes. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
