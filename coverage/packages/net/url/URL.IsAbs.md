# Independent review of `net/url.URL.IsAbs`

**Verdict: FAITHFUL**

Evidence SHA-256: `5ed5999d6f0a7a5686cbfc4f6c21c553593657cdafdda7ac78c011a0a457900e`

Tests only whether the scheme field is nonempty. A nonempty field containing a NUL byte qualifies. A nil receiver panics. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
