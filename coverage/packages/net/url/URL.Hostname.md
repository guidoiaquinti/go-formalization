# Independent review of `net/url.URL.Hostname`

**Verdict: FAITHFUL**

Evidence SHA-256: `a68e74471e02827ed623378df9205d9f88e30ada1aa3ef1e4c78c9606e9cb8c8`

Returns the resulting host and strips matching outer brackets only from that component. The native controls confirmed a bracketed literal with a bad port remains whole and a nil receiver panics. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
