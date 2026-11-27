# Independent review of `net/url.getScheme`

**Verdict: FAITHFUL**

Evidence SHA-256: `3322a1ad6714162062ea9fc63987ed9c4f65d061fe570d6f3aeb9deafc000e68`

A leading colon yields the exact missing-scheme error. The source retains valid scheme case and the raw suffix. A leading ASCII digit keeps the unsplit original path. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
