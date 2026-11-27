# Independent review of `net/url.validOptionalPort`

**Verdict: FAITHFUL**

Evidence SHA-256: `c6dc214d498fb050ec4e9e4f3a844ff00da1fa7aab1134c90d3125672f76d0bd`

Accepts a bare colon and a colon followed by an arbitrarily long ASCII digit sequence without numerical conversion. Both native controls reject an implementation that demands a digit after the colon or parses an integer for range. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
