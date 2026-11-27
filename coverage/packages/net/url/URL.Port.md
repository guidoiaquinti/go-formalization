# Independent review of `net/url.URL.Port`

**Verdict: FAITHFUL**

Evidence SHA-256: `51bb749718b81a44ad0c8fd9de12d33b0ce79dd24cfe90886df3109b1552e60a`

Preserves the exact original digits, including leading zeros. An invalid port leaves the result empty and a nil receiver panics. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
